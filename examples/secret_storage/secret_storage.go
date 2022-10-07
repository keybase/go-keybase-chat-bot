/*
WHAT IS IN THIS EXAMPLE?

Keybase has added an encrypted key-value store intended to support
security-conscious bot development with persistent state. It is a
place to store small bits of data that are

	(1) encrypted for a team or user (via the user's implicit self-team: e.g.

alice,alice),

	(2) persistent across logins
	(3) fast and durable.

It supports putting, getting, listing, and deleting. A team has many
namespaces, a namespace has many entryKeys, and an entryKey has one current
entryValue. Namespaces and entryKeys are in cleartext, and the Keybase client
service will encrypt and sign the entryValue on the way in (as well as
decrypt and verify on the way out) so keybase servers cannot see it or
forge it.

-----------

This example implements a simple bot to manage hackerspace tool rentals. It
shows one way you can obfuscate entryKeys (which are not encrypted) by
storing their HMACs, so that no one but your team (not even
Keybase) can know about the names of all the cool tools you have; you can do
something similar to hide namespaces.

Additionally this example handles concurrent writes by using explicit revision
numbers to prevent one user from unintentionally clobbering another user's
rental updates.

Here we've stored the HMAC secret and other entries in the team's kvstore; you
could also store the entries in the bot's own kvstore (the default team).
*/
package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/keybase1"

	"golang.org/x/sync/errgroup"
)

// SecretKeyKVStoreAPI implements the KVStoreAPI interface, and that hides the
// plaintext of entryKeys from Keybase servers.  It does so by HMACing
// entryKeys using a per-(team, namespace) secret, and storing the HMAC instead
// of the plaintext entryKey. This approach does not handle any secret
// rotation, and does not expect the secret to change. The secret is stored
// under the entryKey "_secret".
//
// The plaintext entryKey is stored in a JSON entryValue under the entryKey
// "_key" to enable listing.
//
// This approach does not hide memory access patterns. Also, Keybase servers
// prevent a removed team member from continuing to access a team's data, but
// if that were somehow bypassed*, a former team member who still knows the
// HMAC secret could check for the presence of specific entryKeys (*but you
// probably have bigger issues to deal with in that case...).
type SecretKeyKVStoreAPI struct {
	api     kbchat.KVStoreAPI
	secrets map[string](map[string][]byte)
	config  SecretKeyKVStoreAPIConfig
}

type SecretKeyKVStoreAPIConfig struct {
	secretName            string
	plaintextEntryKeyName string
}

func NewSecretKeyKVStoreAPI(api kbchat.KVStoreAPI) *SecretKeyKVStoreAPI {
	config := SecretKeyKVStoreAPIConfig{
		secretName:            "_secret",
		plaintextEntryKeyName: "_key",
	}
	sc := SecretKeyKVStoreAPI{
		api:     api,
		secrets: make(map[string](map[string][]byte)),
		config:  config,
	}
	return &sc
}

func (sc *SecretKeyKVStoreAPI) loadSecret(teamName string, namespace string) ([]byte, error) {
	if _, ok := sc.secrets[teamName]; !ok {
		sc.secrets[teamName] = make(map[string][]byte)
	}
	if secret, ok := sc.secrets[teamName][namespace]; ok {
		return secret, nil
	}

	newSecret := make([]byte, sha256.BlockSize)
	if _, err := rand.Read(newSecret); err != nil {
		return nil, err
	}

	// we don't expect SecretKey's revision > 0
	_, err := sc.api.PutEntryWithRevision(&teamName, namespace, sc.config.secretName, hex.EncodeToString(newSecret), 1)
	if err != nil {
		if e, ok := err.(kbchat.Error); !ok || e.Code != kbchat.RevisionErrorCode {
			// unexpected error
			return nil, err
		}

		// failed to put; get entry
		res, err := sc.api.GetEntry(&teamName, namespace, sc.config.secretName)
		if err != nil {
			return nil, err
		}
		var entryValue string
		if res.EntryValue != nil {
			entryValue = *res.EntryValue
		}
		existingSecret, err := hex.DecodeString(entryValue)
		if err != nil {
			return nil, err
		}
		sc.secrets[teamName][namespace] = existingSecret
	} else {
		sc.secrets[teamName][namespace] = newSecret
	}
	return sc.secrets[teamName][namespace], nil
}

func (sc *SecretKeyKVStoreAPI) obfuscateEntryKey(teamName string, namespace string, entryKey string) (string, error) {
	secret, err := sc.loadSecret(teamName, namespace)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(entryKey))
	hmacEntryKey := mac.Sum(nil)
	return hex.EncodeToString(hmacEntryKey), nil
}

func (sc *SecretKeyKVStoreAPI) PutEntry(teamName *string, namespace string, entryKey string, entryValue string) (result keybase1.KVPutResult, err error) {
	return sc.PutEntryWithRevision(teamName, namespace, entryKey, entryValue, 0)
}

func (sc *SecretKeyKVStoreAPI) PutEntryWithRevision(teamName *string, namespace string, entryKey string, entryValue string, revision int) (result keybase1.KVPutResult, err error) {
	if teamName == nil {
		return result, fmt.Errorf("teamName must be defined for SecretKeyKVStoreAPI methods")
	}
	return sc.PutEntryWithRevisionAndTeam(*teamName, namespace, entryKey, entryValue, revision)
}

func (sc *SecretKeyKVStoreAPI) PutEntryWithRevisionAndTeam(teamName string, namespace string, entryKey string, entryValue string, revision int) (result keybase1.KVPutResult, err error) {
	var keyedValue map[string]string
	if err = json.Unmarshal([]byte(entryValue), &keyedValue); err != nil {
		return result, err
	}

	keyedValue[sc.config.plaintextEntryKeyName] = entryKey
	bytes, err := json.Marshal(keyedValue)
	if err != nil {
		return result, err
	}
	hmacEntryKey, err := sc.obfuscateEntryKey(teamName, namespace, entryKey)
	if err != nil {
		return result, err
	}

	result, err = sc.api.PutEntryWithRevision(&teamName, namespace, hmacEntryKey, string(bytes), revision)
	if err != nil {
		return result, err
	}
	result.EntryKey = entryKey
	return result, err
}

func (sc *SecretKeyKVStoreAPI) DeleteEntry(teamName *string, namespace string, entryKey string) (result keybase1.KVDeleteEntryResult, err error) {
	return sc.DeleteEntryWithRevision(teamName, namespace, entryKey, 0)
}

func (sc *SecretKeyKVStoreAPI) DeleteEntryWithRevision(teamName *string, namespace string, entryKey string, revision int) (result keybase1.KVDeleteEntryResult, err error) {
	if teamName == nil {
		return result, fmt.Errorf("teamName must be defined for SecretKeyKVStoreAPI methods")
	}
	return sc.DeleteEntryWithRevisionAndTeam(*teamName, namespace, entryKey, revision)
}

func (sc *SecretKeyKVStoreAPI) DeleteEntryWithRevisionAndTeam(teamName string, namespace string, entryKey string, revision int) (result keybase1.KVDeleteEntryResult, err error) {
	hmacEntryKey, err := sc.obfuscateEntryKey(teamName, namespace, entryKey)
	if err != nil {
		return result, err
	}
	result, err = sc.api.DeleteEntryWithRevision(&teamName, namespace, hmacEntryKey, revision)
	if err != nil {
		return result, err
	}
	result.EntryKey = entryKey
	return result, err
}

func (sc *SecretKeyKVStoreAPI) GetEntry(teamName *string, namespace string, entryKey string) (result keybase1.KVGetResult, err error) {
	if teamName == nil {
		return result, fmt.Errorf("teamName must be defined for SecretKeyKVStoreAPI methods")
	}
	return sc.GetEntryWithTeam(*teamName, namespace, entryKey)
}

func (sc *SecretKeyKVStoreAPI) GetEntryWithTeam(teamName string, namespace string, entryKey string) (result keybase1.KVGetResult, err error) {
	hmacEntryKey, err := sc.obfuscateEntryKey(teamName, namespace, entryKey)
	if err != nil {
		return result, err
	}
	result, err = sc.api.GetEntry(&teamName, namespace, hmacEntryKey)
	if err != nil {
		return result, err
	}
	result.EntryKey = entryKey
	return result, err
}

func (sc *SecretKeyKVStoreAPI) ListNamespaces(teamName *string) (keybase1.KVListNamespaceResult, error) {
	return sc.ListNamespaces(teamName)
}

func (sc *SecretKeyKVStoreAPI) ListEntryKeys(teamName *string, namespace string) (result keybase1.KVListEntryResult, err error) {
	keys, err := sc.api.ListEntryKeys(teamName, namespace)
	if err != nil {
		return result, err
	}
	tmp := keys.EntryKeys[:0]
	for _, e := range keys.EntryKeys {
		if strings.HasPrefix(e.EntryKey, "_") {
			continue
		}

		get, err := sc.api.GetEntry(teamName, namespace, e.EntryKey)
		if err != nil {
			return result, err
		}
		var entryValue string
		if get.EntryValue != nil {
			entryValue = *get.EntryValue
		}
		var keyedValue map[string]string
		if err := json.Unmarshal([]byte(entryValue), &keyedValue); err != nil {
			return result, err
		}
		e.EntryKey = keyedValue[sc.config.plaintextEntryKeyName]
		tmp = append(tmp, e)
	}
	keys.EntryKeys = tmp
	return keys, nil
}

// RentalBotClient wraps a KVStoreClient to expose methods to handle tool rentals.
// Tries kvstore write actions with explicit revision numbers.
// If it fails to write, it does a "get" and returns the get result.
type RentalBotClient struct {
	api       kbchat.KVStoreAPI
	team      string
	namespace string
}

func NewRentalBotClient(api kbchat.KVStoreAPI, teamName string, namespace string) *RentalBotClient {
	r := RentalBotClient{api: api, team: teamName, namespace: namespace}
	return &r
}

func (r *RentalBotClient) Lookup(tool string) (keybase1.KVGetResult, error) {
	return r.api.GetEntry(&r.team, r.namespace, tool)
}

// Add returns (whether action is successful, most recent get result if applicable, error)
func (r *RentalBotClient) Add(tool string) (ok bool, result keybase1.KVGetResult, err error) {
	result, err = r.Lookup(tool)
	if err != nil {
		return false, result, err // api call failed
	} else if result.EntryValue != nil {
		return true, result, nil // tool already exists
	}

	expectedRevision := result.Revision + 1
	val := make(map[string]string)
	bytes, err := json.Marshal(val)
	if err != nil {
		return false, result, err
	}
	_, err = r.api.PutEntryWithRevision(&r.team, r.namespace, tool, string(bytes), expectedRevision)
	if err != nil {
		if e, ok := err.(kbchat.Error); !ok || e.Code != kbchat.RevisionErrorCode {
			// unexpected error
			return false, result, err
		}

		// failed put. try get
		result, err := r.Lookup(tool)
		if err != nil {
			return false, result, err // api call failed
		}
		return false, result, nil // failed put. return KVGetResult
	}
	return true, result, nil // successul put
}

// Remove returns (whether action is successful, most recent get result if applicable, error)
func (r *RentalBotClient) Remove(tool string) (ok bool, result keybase1.KVGetResult, err error) {
	result, err = r.Lookup(tool)
	if err != nil {
		return false, result, err // api call failed
	} else if result.EntryValue == nil {
		return true, result, nil // tool already doesn't exist
	}

	expectedRevision := result.Revision + 1

	_, err = r.api.DeleteEntryWithRevision(&r.team, r.namespace, tool, expectedRevision)
	switch err.(type) {
	case nil:
		// successul delete
		return true, result, nil
	case kbchat.Error:
		// failed delete. try get
		result, err := r.Lookup(tool)
		if err != nil {
			return false, result, err // api call failed
		}
		return false, result, nil // failed delete. return KVGetResult
	default:
		// unexpected error
		return false, result, err
	}
}

// Reserve reserves a tool for a given day if that day is not already reserved.
// Note: if you reserve a not-added or deleted tool, it will add the tool.
// Returns (whether action is successful, most recent get result if applicable, error)
func (r *RentalBotClient) Reserve(username string, tool string, day string) (ok bool, result keybase1.KVGetResult, err error) {
	var reservations map[string]string
	result, err = r.Lookup(tool)
	if err != nil {
		return false, result, err // api call failed
	}
	var entryValue string
	if result.EntryValue != nil {
		entryValue = *result.EntryValue
	}
	reservations = make(map[string]string)
	if err = json.Unmarshal([]byte(entryValue), &reservations); err != nil {
		return false, result, err
	}

	if _, ok := reservations[day]; ok {
		return false, result, nil // failed to put because day is already reserved
	}
	reservations[day] = username
	expectedRevision := result.Revision + 1

	bytes, err := json.Marshal(reservations)
	if err != nil {
		return false, result, err
	}
	_, err = r.api.PutEntryWithRevision(&r.team, r.namespace, tool, string(bytes), expectedRevision)
	if err != nil {
		if e, ok := err.(kbchat.Error); !ok || e.Code != kbchat.RevisionErrorCode {
			// unexpected error
			return false, result, err
		}

		// failed put. try get
		result, err := r.Lookup(tool)
		if err != nil {
			return false, result, err // api call failed
		}
		return false, result, nil // failed put. return KVGetResult
	}
	return true, result, nil // successul put
}

// Unreserve a tool for a given day if that day is currently reserved by the given user.
// Note: if you unreserve a not-added or deleted tool, it will not add the tool.
// Returns (whether action is successful, most recent get result if applicable, error)
func (r *RentalBotClient) Unreserve(username string, tool string, day string) (ok bool, result keybase1.KVGetResult, err error) {
	var reservations map[string]string
	result, err = r.Lookup(tool)
	if err != nil {
		return false, result, err // api call failed
	}
	var entryValue string
	if result.EntryValue != nil {
		entryValue = *result.EntryValue
	}
	reservations = make(map[string]string)
	if err = json.Unmarshal([]byte(entryValue), &reservations); err != nil {
		return false, result, err
	}

	reserver, ok := reservations[day]
	if !ok {
		// a noop because currently not reserved
		return true, result, nil
	} else if reserver != username {
		// failed to put because current reserver is not user
		return false, result, nil
	}
	expectedRevision := result.Revision + 1
	delete(reservations, day)

	bytes, err := json.Marshal(reservations)
	if err != nil {
		return false, result, err
	}

	_, err = r.api.PutEntryWithRevision(&r.team, r.namespace, tool, string(bytes), expectedRevision)
	if err != nil {
		if e, ok := err.(kbchat.Error); !ok || e.Code != kbchat.RevisionErrorCode {
			// unexpected error
			return false, result, err
		}

		// failed put. try get
		result, err = r.Lookup(tool)
		if err != nil {
			return false, result, err // api call failed
		}
		return false, result, nil // failed put. return KVGetResult
	}
	return true, result, nil // successul put
}

func (r *RentalBotClient) ListTools() ([]string, error) {
	var tools []string
	res, err := r.api.ListEntryKeys(&r.team, r.namespace)

	if err != nil {
		return tools, err
	}

	for _, tool := range res.EntryKeys {
		tools = append(tools, tool.EntryKey)
	}
	return tools, nil
}

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func basicRentalUsers(rental *RentalBotClient) error {
	user1 := "Jo"
	user2 := "Charlie"
	date1 := "2044-03-12"
	date2 := "2044-06-12"
	date3 := "2044-06-13"
	tool := "laz0rs"

	ok, res, err := rental.Remove(tool)
	fmt.Printf("REMOVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	var tools []string
	tools, err = rental.ListTools()
	fmt.Printf("LIST TOOLS: %+v, %v\n", tools, err)
	if err != nil {
		return fmt.Errorf("Unexpected result: %v", err)
	}

	res, err = rental.Lookup(tool)
	fmt.Printf("LOOKUP: %+v, %v\n", res, err)
	if err != nil {
		return fmt.Errorf("Unexpected result: %v", err)
	}

	ok, res, err = rental.Add("time travel machine")
	fmt.Printf("ADD: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Remove(tool)
	fmt.Printf("REMOVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Add(tool)
	fmt.Printf("ADD: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Reserve(user1, tool, date1)
	fmt.Printf("RESERVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Reserve(user1, tool, date1)
	fmt.Printf("EXPECTING RESERVE FAIL: %v, %+v, %v\n", ok, res, err)
	if ok && err == nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Reserve(user2, tool, date2)
	fmt.Printf("RESERVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	res, err = rental.Lookup(tool)
	fmt.Printf("LOOKUP: %+v, %v\n", res, err)

	ok, res, err = rental.Unreserve(user1, tool, date3)
	fmt.Printf("UNRESERVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Unreserve(user1, tool, date2)
	fmt.Printf("EXPECTING UNRESERVE FAIL: %v, %+v, %v\n", ok, res, err)
	if ok && err == nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Unreserve(user1, tool, date1)
	fmt.Printf("UNRESERVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	res, err = rental.Lookup(tool)
	fmt.Printf("LOOKUP: %+v, %v\n", res, err)
	return nil
}

func concurrentRentalUsers(rental *RentalBotClient) error {
	tool := "time travel machine"
	var g errgroup.Group

	// pre
	for {
		ok, res, err := rental.Remove(tool)
		fmt.Printf("TRY TO REMOVE: %v, %+v, %+v\n", ok, res, err)
		if ok && err == nil {
			break
		}
	}

	// have 5 users concurrently try to reserve the same tool for 5 unique dates
	for i := 1; i <= 5; i++ {
		g.Go(func(userID int) func() error {
			return func() error {
				date := fmt.Sprintf("2044-10-0%d", userID)
				user := fmt.Sprintf("user%d", userID)

				i := 0
				// keep trying to reserve for user's unique date until successful
				for {
					ok, res, err := rental.Reserve(user, tool, date)
					i++
					fmt.Printf("%v, attempt %d, TRY TO RESERVE: %v, %+v, %+v\n", user, i, ok, res, err)
					if ok && err == nil {
						break
					}
				}
				return nil
			}
		}(i))
	}
	g.Wait()

	// post: check that the tool has been reserved for all 5 unique dates
	var val map[string]string
	res, err := rental.Lookup(tool)
	if err != nil {
		return err
	} else if res.EntryValue == nil {
		val = make(map[string]string)
	} else {
		var entryValue string
		if res.EntryValue != nil {
			entryValue = *res.EntryValue
		}
		if err := json.Unmarshal([]byte(entryValue), &val); err != nil {
			return err
		}
	}
	if len(val) != 6 {
		return fmt.Errorf("Unexpected result: %+v", val)
	}
	return nil
}

func main() {
	const MsgPrefix = "!storage"

	var kbLoc string
	var kbc *kbchat.API
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	fmt.Println("Starting secret_storage example...")

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	team := "yourhackerspace"

	secretClient := NewSecretKeyKVStoreAPI(kbc)
	rental := NewRentalBotClient(secretClient, team, "!rental")

	fmt.Println("...basic rental actions...")
	if err = basicRentalUsers(rental); err != nil {
		log.Fatalf("Bot failed: %+v", err)
	}
	fmt.Println("...multiple users try to reserve...")
	if err = concurrentRentalUsers(rental); err != nil {
		log.Fatalf("Bot failed: %+v", err)
	}
	fmt.Println("...secret_storage example is complete.")
}
