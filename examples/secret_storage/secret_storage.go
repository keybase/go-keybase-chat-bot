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

// SecretKeyKVStoreAPI is a KVStoreAPI that hides the entryKeys from Keybase servers.
// It does so by HMACing entryKeys using a per-(team, namespace) secret,
// and storing the HMAC instead of the plaintext entryKey. This approach
// does not handle any secret rotation, and does not expect the secret to
// change.
//
// The plaintext entryKey is stored in it's corresponding JSON entryValue
// under the key "_key" to enable listing.
//
// This approach does not hide memory access patterns. Also, Keybase
// servers prevent a removed team member from continuing to access a team's
// data, but if that were somehow bypassed*, a former team member who still
// knows the HMAC secret can check for the presence of specific entryKeys
// (*but you probably have bigger issues to deal with in that case...).
type SecretKeyKVStoreAPI struct {
	api       kbchat.KVStoreAPI
	secrets   map[string](map[string][]byte)
	secretKey string
	keyKey    string
}

func NewSecretKeyKVStoreAPI(api kbchat.KVStoreAPI) *SecretKeyKVStoreAPI {
	secrets := make(map[string](map[string][]byte))
	sc := SecretKeyKVStoreAPI{api, secrets, "_secret", "_key"}
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
	_, err := sc.api.PutEntryWithRevision(teamName, namespace, sc.secretKey, hex.EncodeToString(newSecret), 1)
	if err != nil {
		if e, ok := err.(kbchat.Error); !ok || e.Code != kbchat.RevisionErrorCode {
			// unexpected error
			return nil, err
		}

		// failed to put; get entry
		res, err := sc.api.GetEntry(teamName, namespace, sc.secretKey)
		if err != nil {
			return nil, err
		}
		existingSecret, err := hex.DecodeString(res.EntryValue)
		if err != nil {
			return nil, err
		}
		sc.secrets[teamName][namespace] = existingSecret
	} else {
		sc.secrets[teamName][namespace] = newSecret
	}
	return sc.secrets[teamName][namespace], nil
}

func (sc *SecretKeyKVStoreAPI) hmacKey(teamName string, namespace string, entryKey string) (string, error) {
	secret, err := sc.loadSecret(teamName, namespace)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(entryKey))
	hmacEntryKey := mac.Sum(nil)
	return hex.EncodeToString(hmacEntryKey), nil
}

func (sc *SecretKeyKVStoreAPI) PutEntry(teamName string, namespace string, entryKey string, entryValue string) (result keybase1.KVPutResult, err error) {
	return sc.PutEntryWithRevision(teamName, namespace, entryKey, entryValue, 0)
}

func (sc *SecretKeyKVStoreAPI) PutEntryWithRevision(teamName string, namespace string, entryKey string, entryValue string, revision int) (result keybase1.KVPutResult, err error) {
	var keyedValue map[string]string
	if err = json.Unmarshal([]byte(entryValue), &keyedValue); err != nil {
		return result, err
	}

	keyedValue[sc.keyKey] = entryKey
	bytes, err := json.Marshal(keyedValue)
	if err != nil {
		return result, err
	}

	hmacEntryKey, err := sc.hmacKey(teamName, namespace, entryKey)
	if err != nil {
		return result, err
	}

	result, err = sc.api.PutEntryWithRevision(teamName, namespace, hmacEntryKey, string(bytes), revision)
	if err != nil {
		return result, err
	}
	result.EntryKey = entryKey
	return result, err
}

func (sc *SecretKeyKVStoreAPI) DeleteEntry(teamName string, namespace string, entryKey string) (result keybase1.KVDeleteEntryResult, err error) {
	return sc.DeleteEntryWithRevision(teamName, namespace, entryKey, 0)
}

func (sc *SecretKeyKVStoreAPI) DeleteEntryWithRevision(teamName string, namespace string, entryKey string, revision int) (result keybase1.KVDeleteEntryResult, err error) {
	hmacEntryKey, err := sc.hmacKey(teamName, namespace, entryKey)
	if err != nil {
		return result, err
	}
	result, err = sc.api.DeleteEntryWithRevision(teamName, namespace, hmacEntryKey, revision)
	if err != nil {
		return result, err
	}
	result.EntryKey = entryKey
	return result, err
}

func (sc *SecretKeyKVStoreAPI) GetEntry(teamName string, namespace string, entryKey string) (result keybase1.KVGetResult, err error) {
	hmacEntryKey, err := sc.hmacKey(teamName, namespace, entryKey)
	if err != nil {
		return result, err
	}
	result, err = sc.api.GetEntry(teamName, namespace, hmacEntryKey)
	if err != nil {
		return result, err
	}
	result.EntryKey = entryKey
	return result, err
}

func (sc *SecretKeyKVStoreAPI) ListNamespaces(teamName string) (keybase1.KVListNamespaceResult, error) {
	return sc.ListNamespaces(teamName)
}

func (sc *SecretKeyKVStoreAPI) ListEntryKeys(teamName string, namespace string) (result keybase1.KVListEntryResult, err error) {
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

		var keyedValue map[string]string
		if err := json.Unmarshal([]byte(get.EntryValue), &keyedValue); err != nil {
			return result, err
		}
		e.EntryKey = keyedValue[sc.keyKey]
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
	namespace string
}

func NewRentalBotClient(api kbchat.KVStoreAPI, namespace string) *RentalBotClient {
	r := RentalBotClient{api, namespace}
	return &r
}

func (r *RentalBotClient) Lookup(teamName string, tool string) (keybase1.KVGetResult, error) {
	return r.api.GetEntry(teamName, r.namespace, tool)
}

// Add returns (whether action is successful, most recent get result if applicable, error)
func (r *RentalBotClient) Add(teamName string, tool string) (ok bool, result keybase1.KVGetResult, err error) {
	result, err = r.Lookup(teamName, tool)
	if err != nil {
		return false, result, err // api call failed
	} else if result.EntryValue != "" {
		return true, result, nil // tool already exists
	}

	expectedRevision := result.Revision + 1
	val := make(map[string]string)
	bytes, err := json.Marshal(val)
	if err != nil {
		return false, result, err
	}
	_, err = r.api.PutEntryWithRevision(teamName, r.namespace, tool, string(bytes), expectedRevision)
	if err != nil {
		if e, ok := err.(kbchat.Error); !ok || e.Code != kbchat.RevisionErrorCode {
			// unexpected error
			return false, result, err
		}

		// failed put. try get
		result, err := r.Lookup(teamName, tool)
		if err != nil {
			return false, result, err // api call failed
		}
		return false, result, nil // failed put. return KVGetResult
	}
	return true, result, nil // successul put
}

// Remove returns (whether action is successful, most recent get result if applicable, error)
func (r *RentalBotClient) Remove(teamName string, tool string) (ok bool, result keybase1.KVGetResult, err error) {
	result, err = r.Lookup(teamName, tool)
	if err != nil {
		return false, result, err // api call failed
	} else if result.EntryValue == "" {
		return true, result, nil // tool already doesn't exist
	}

	expectedRevision := result.Revision + 1
	_, err = r.api.DeleteEntryWithRevision(teamName, r.namespace, tool, expectedRevision)
	if err != nil {
		if _, ok := err.(kbchat.Error); !ok {
			// unexpected error
			return false, result, err
		}

		// failed delete. try get
		result, err := r.Lookup(teamName, tool)
		if err != nil {
			return false, result, err // api call failed
		}
		return false, result, nil // failed delete. return KVGetResult
	}
	return true, result, nil // successul delete
}

// Reserve reserve a tool for a given day if that day is already not reserved.
// Note: if you reserve a not-added or deleted tool, it will add the tool.
// Returns (whether action is successful, most recent get result if applicable, error)
func (r *RentalBotClient) Reserve(teamName string, username string, tool string, day string) (ok bool, result keybase1.KVGetResult, err error) {
	var val map[string]string
	result, err = r.Lookup(teamName, tool)
	if err != nil {
		return false, result, err // api call failed
	} else if result.EntryValue == "" {
		val = make(map[string]string)
	} else {
		if err = json.Unmarshal([]byte(result.EntryValue), &val); err != nil {
			return false, result, err
		}
	}

	if _, ok := val[day]; ok {
		return false, result, nil // failed to put because day is already reserved
	}
	val[day] = username
	expectedRevision := result.Revision + 1

	bytes, err := json.Marshal(val)
	if err != nil {
		return false, result, err
	}
	_, err = r.api.PutEntryWithRevision(teamName, r.namespace, tool, string(bytes), expectedRevision)
	if err != nil {
		if e, ok := err.(kbchat.Error); !ok || e.Code != kbchat.RevisionErrorCode {
			// unexpected error
			return false, result, err
		}

		// failed put. try get
		result, err := r.Lookup(teamName, tool)
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
func (r *RentalBotClient) Unreserve(teamName string, username string, tool string, day string) (ok bool, result keybase1.KVGetResult, err error) {
	var val map[string]string
	result, err = r.Lookup(teamName, tool)
	if err != nil {
		return false, result, err // api call failed
	} else if result.EntryValue == "" {
		val = make(map[string]string)
	} else {
		if err := json.Unmarshal([]byte(result.EntryValue), &val); err != nil {
			return false, result, err
		}
	}

	reserver, ok := val[day]
	if !ok {
		// a noop because currently not reserved
		return true, result, nil
	} else if reserver != username {
		// failed to put because current reserver is not user
		return false, result, nil
	}
	expectedRevision := result.Revision + 1
	delete(val, day)

	bytes, err := json.Marshal(val)
	if err != nil {
		return false, result, err
	}

	_, err = r.api.PutEntryWithRevision(teamName, r.namespace, tool, string(bytes), expectedRevision)
	if err != nil {
		if e, ok := err.(kbchat.Error); !ok || e.Code != kbchat.RevisionErrorCode {
			// unexpected error
			return false, result, err
		}

		// failed put. try get
		result, err = r.Lookup(teamName, tool)
		if err != nil {
			return false, result, err // api call failed
		}
		return false, result, nil // failed put. return KVGetResult
	}
	return true, result, nil // successul put
}

func (r *RentalBotClient) ListTools(teamName string) ([]string, error) {
	var tools []string
	res, err := r.api.ListEntryKeys(teamName, r.namespace)

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

func basicRentalUsers(rental *RentalBotClient, team string) error {
	user1 := "Jo"
	user2 := "Charlie"
	date1 := "2044-03-12"
	date2 := "2044-06-12"
	date3 := "2044-06-13"
	tool := "laz0rs"

	ok, res, err := rental.Remove(team, tool)
	fmt.Printf("REMOVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	var tools []string
	tools, err = rental.ListTools(team)
	fmt.Printf("LIST TOOLS: %+v, %v\n", tools, err)

	res, err = rental.Lookup(team, tool)
	fmt.Printf("LOOKUP: %+v, %v\n", res, err)

	ok, res, err = rental.Add(team, "time travel machine")
	fmt.Printf("ADD: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Remove(team, tool)
	fmt.Printf("REMOVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Add(team, tool)
	fmt.Printf("ADD: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Reserve(team, user1, tool, date1)
	fmt.Printf("RESERVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Reserve(team, user1, tool, date1)
	fmt.Printf("EXPECTING RESERVE FAIL: %v, %+v, %v\n", ok, res, err)
	if ok && err == nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Reserve(team, user2, tool, date2)
	fmt.Printf("RESERVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	res, err = rental.Lookup(team, tool)
	fmt.Printf("LOOKUP: %+v, %v\n", res, err)

	ok, res, err = rental.Unreserve(team, user1, tool, date3)
	fmt.Printf("UNRESERVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Unreserve(team, user1, tool, date2)
	fmt.Printf("EXPECTING UNRESERVE FAIL: %v, %+v, %v\n", ok, res, err)
	if ok && err == nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	ok, res, err = rental.Unreserve(team, user1, tool, date1)
	fmt.Printf("UNRESERVE: %v, %+v, %v\n", ok, res, err)
	if !ok || err != nil {
		return fmt.Errorf("Unexpected result: %v, %v", ok, err)
	}

	res, err = rental.Lookup(team, tool)
	fmt.Printf("LOOKUP: %+v, %v\n", res, err)
	return nil
}

func concurrentRentalUsers(rental *RentalBotClient, team string) error {
	tool := "time travel machine"
	var g errgroup.Group

	// pre
	for {
		ok, res, err := rental.Remove(team, tool)
		fmt.Printf("TRY TO REMOVE: %v, %+v, %+v\n", ok, res, err)
		if ok && err == nil {
			break
		}
	}
	g.Wait()

	// have 5 users concurrently try to reserve the same tool for 5 unique dates
	for i := 1; i <= 5; i++ {
		g.Go(func(userID int) func() error {
			return func() error {
				date := fmt.Sprintf("2044-10-0%d", userID)
				user := fmt.Sprintf("user%d", userID)

				i := 0
				// keep trying to reserve for user's unique date until successful
				for {
					ok, res, err := rental.Reserve(team, user, tool, date)
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
	res, err := rental.Lookup(team, tool)
	if err != nil {
		return err
	} else if res.EntryValue == "" {
		val = make(map[string]string)
	} else {
		if err := json.Unmarshal([]byte(res.EntryValue), &val); err != nil {
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
	rental := NewRentalBotClient(secretClient, "!rental")

	fmt.Println("...basic rental actions...")
	if err = basicRentalUsers(rental, team); err != nil {
		log.Fatalf("Bot failed: %+v", err)
	}
	fmt.Println("...multiple users try to reserve...")
	if err = concurrentRentalUsers(rental, team); err != nil {
		log.Fatalf("Bot failed: %+v", err)
	}
	fmt.Println("...secret_storage example is complete.")
}
