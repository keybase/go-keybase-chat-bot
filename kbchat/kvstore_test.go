package kbchat

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/keybase1"
	"github.com/stretchr/testify/require"
)

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func containsKey(a []keybase1.KVListEntryKey, x string) bool {
	for _, n := range a {
		if x == n.EntryKey {
			return true
		}
	}
	return false
}

func clearNamespace(bot *API, team string, namespace string) error {
	fmt.Printf("Clearing namespace %s\n", namespace)
	listek, err := bot.ListEntryKeys(team, namespace)

	for _, entryKey := range listek.EntryKeys {
		if _, err = bot.DeleteEntry(team, namespace, entryKey.EntryKey); err != nil {
			return err
		}
	}
	return nil
}

// test doesn't make assumption about prior state
func TestKVStore(t *testing.T) {
	var err error

	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getTeamChatChannel(t, "acme")
	team := channel.Name

	namespace := fmt.Sprintf("_test_namespace%d", rand.Int())
	entryKey := "_test_key1"

	require.NoError(t, clearNamespace(alice, team, namespace))
	defer func() {
		require.NoError(t, clearNamespace(alice, team, namespace))
	}()

	// put with default revision
	put, err := alice.PutEntry(team, namespace, entryKey, "value1")
	require.NoError(t, err)
	require.True(t, put.Revision > 0)
	currentRevision := put.Revision

	expectedRevision := currentRevision + 1

	// fail put (wrong revision)
	_, err = alice.PutEntryWithRevision(team, namespace, entryKey, "value2", expectedRevision-1)
	require.Error(t, err)
	require.Equal(t, RevisionErrorCode, err.(Error).Code)

	// list namespaces
	listns, err := alice.ListNamespaces(team)
	require.NoError(t, err)
	require.True(t, len(listns.Namespaces) > 0)
	require.True(t, contains(listns.Namespaces, namespace))

	// list entryKeys
	listek, err := alice.ListEntryKeys(team, namespace)
	require.NoError(t, err)
	require.True(t, len(listek.EntryKeys) > 0)
	require.True(t, containsKey(listek.EntryKeys, entryKey))

	// get
	get, err := alice.GetEntry(team, namespace, entryKey)
	require.NoError(t, err)
	require.Equal(t, "value1", get.EntryValue)

	// fail delete
	_, err = alice.DeleteEntryWithRevision(team, namespace, entryKey, expectedRevision+1)
	require.Error(t, err)
	require.Equal(t, RevisionErrorCode, err.(Error).Code)

	// delete
	del, err := alice.DeleteEntryWithRevision(team, namespace, entryKey, expectedRevision)
	require.NoError(t, err)
	require.Equal(t, expectedRevision, del.Revision)

	// fail delete (non existent)
	_, err = alice.DeleteEntry(team, namespace, entryKey)
	require.Error(t, err)
	require.Equal(t, DeleteNonExistentErrorCode, err.(Error).Code)

	// put with default revision
	expectedRevision++
	put, err = alice.PutEntry(team, namespace, entryKey, "value3")
	require.NoError(t, err)
	require.Equal(t, expectedRevision, put.Revision)

	// delete with default revision
	expectedRevision++
	del, err = alice.DeleteEntry(team, namespace, entryKey)
	require.NoError(t, err)
	require.Equal(t, expectedRevision, del.Revision)

	// get
	get, err = alice.GetEntry(team, namespace, entryKey)
	require.NoError(t, err)
	require.Equal(t, "", get.EntryValue)
	require.Equal(t, expectedRevision, get.Revision)
}
