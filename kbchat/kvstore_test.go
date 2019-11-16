package kbchat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// test doesn't make assumption about prior state
func TestKVStore(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	channel := getTeamChatChannel(t, "acme")
	team := channel.Name
	defer testBotTeardown(t, alice, dir)

	namespace := "_test_namespace1"
	key := "_test_key1"

	// put with default revision
	res2, err2 := alice.PutEntry(team, namespace, key, "value1", 0)
	require.NoError(t, err2)
	rev := res2.Revision

	// fail put
	_, err3 := alice.PutEntry(team, namespace, key, "value2", rev)
	require.Error(t, err3)

	// list namespaces
	res4, err4 := alice.ListNamespaces(team)
	require.NoError(t, err4)
	require.True(t, len(res4.Namespaces) > 0)

	// list entryKeys
	res5, err5 := alice.ListEntryKeys(team, namespace)
	require.NoError(t, err5)
	require.True(t, len(res5.EntryKeys) > 0)

	// get
	res6, err6 := alice.GetEntry(team, namespace, key)
	require.NoError(t, err6)
	require.Equal(t, "value1", res6.EntryValue)

	// fail delete
	_, err7 := alice.DeleteEntry(team, namespace, key, rev+2)
	require.Error(t, err7)

	// delete
	res8, err8 := alice.DeleteEntry(team, namespace, key, rev+1)
	require.NoError(t, err8)
	require.Equal(t, rev+1, res8.Revision)

	// fail delete
	_, err9 := alice.DeleteEntry(team, namespace, key, 0)
	require.Error(t, err9)

	// get
	res10, err10 := alice.GetEntry(team, namespace, key)
	require.NoError(t, err10)
	require.Equal(t, "", res10.EntryValue)
}
