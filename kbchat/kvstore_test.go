package kbchat

import (
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

// test doesn't make assumption about prior state
func TestKVStore(t *testing.T) {
	var err error

	alice, dir := testBotSetup(t, "alice")
	channel := getTeamChatChannel(t, "acme")
	team := channel.Name
	defer testBotTeardown(t, alice, dir)

	namespace := "_test_namespace1"
	key := "_test_key1"

	// put with default revision
	put, err := alice.PutEntry(team, namespace, key, "value1")
	require.NoError(t, err)
	require.True(t, put.Revision > 0)
	rev := put.Revision

	expectedRevision := rev + 1

	// fail put (wrong revision)
	_, err = alice.PutEntryWithRevision(team, namespace, key, "value2", expectedRevision-1)
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
	require.True(t, containsKey(listek.EntryKeys, key))

	// get
	get, err := alice.GetEntry(team, namespace, key)
	require.NoError(t, err)
	require.Equal(t, "value1", get.EntryValue)

	// fail delete
	_, err = alice.DeleteEntryWithRevision(team, namespace, key, expectedRevision+1)
	require.Error(t, err)
	require.Equal(t, RevisionErrorCode, err.(Error).Code)

	// delete
	del, err := alice.DeleteEntryWithRevision(team, namespace, key, expectedRevision)
	require.NoError(t, err)
	require.Equal(t, expectedRevision, del.Revision)

	// fail delete (non existent)
	_, err = alice.DeleteEntry(team, namespace, key)
	require.Error(t, err)
	require.Equal(t, DeleteNonExistentErrorCode, err.(Error).Code)

	// put with default revision
	expectedRevision++
	put, err = alice.PutEntry(team, namespace, key, "value3")
	require.NoError(t, err)
	require.Equal(t, expectedRevision, put.Revision)

	// delete with default revision
	expectedRevision++
	del, err = alice.DeleteEntry(team, namespace, key)
	require.NoError(t, err)
	require.Equal(t, expectedRevision, del.Revision)

	// get
	get, err = alice.GetEntry(team, namespace, key)
	require.NoError(t, err)
	require.Equal(t, "", get.EntryValue)
	require.Equal(t, expectedRevision, get.Revision)
}
