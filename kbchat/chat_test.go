package kbchat

import (
	"os"
	"path"
	"testing"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
	"github.com/stretchr/testify/require"
)

func TestGetConversations(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	conversations, err := alice.GetConversations(false)
	require.NoError(t, err)
	require.True(t, len(conversations) > 0)
}

func TestGetTextMessages(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getOneOnOneChatChannel(t, "alice", "bob")
	messages, err := alice.GetTextMessages(channel, false)
	require.NoError(t, err)
	require.True(t, len(messages) > 0)
}

func TestSendMessage(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getOneOnOneChatChannel(t, "alice", "bob")
	text := "test SendMessage " + randomString(t)

	// Send the message
	res, err := alice.SendMessage(channel, "%s", text)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(t, alice, channel)
	require.Equal(t, readMessage.Content.TypeName, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.Id, *res.Result.MessageID)
}
func TestSendMessageByConvID(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	text := "test SendMessageByConvID " + randomString(t)
	channel := getOneOnOneChatChannel(t, "alice", "bob")
	convID := getConvIDForChannel(t, alice, channel)

	// Send the message
	res, err := alice.SendMessageByConvID(convID, "%s", text)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(t, alice, channel)
	require.Equal(t, readMessage.Content.TypeName, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.Id, *res.Result.MessageID)
}

func TestSendMessageByTlfName(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	text := "test SendMessageByTlfName " + randomString(t)
	channel := getOneOnOneChatChannel(t, "alice", "bob")

	// Send the message
	res, err := alice.SendMessageByTlfName(channel.Name, "%s", text)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(t, alice, channel)
	require.Equal(t, readMessage.Content.TypeName, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.Id, *res.Result.MessageID)
}

func TestSendMessageByTeamName(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	text := "test SendMessageByTeamName " + randomString(t)
	channel := getTeamChatChannel(t, "acme")

	// Send the message
	res, err := alice.SendMessageByTeamName(channel.Name, &channel.TopicName, "%s", text)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(t, alice, channel)
	require.Equal(t, readMessage.Content.TypeName, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.Id, *res.Result.MessageID)
}

func TestSendAttachmentByTeam(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getTeamChatChannel(t, "acme")

	// Create a test file
	fileName := "kb-attachment.txt"
	location := path.Join(os.TempDir(), fileName)
	data := []byte("My super cool attachment" + randomString(t))
	err := os.WriteFile(location, data, 0644)
	require.NoError(t, err)

	// Send the message
	title := "test SendAttachmentByTeam " + randomString(t)
	res, err := alice.SendAttachmentByTeam(channel.Name, &channel.TopicName, location, title)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)
}

// //////////////////////////////////////////////////////
// React to chat ///////////////////////////////////////
// //////////////////////////////////////////////////////
func TestReactByChannel(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getOneOnOneChatChannel(t, "alice", "bob")

	react := ":cool:"
	lastMessageID := getMostRecentMessage(t, alice, channel).Id

	res, err := alice.ReactByChannel(channel, lastMessageID, react)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)
}

func TestReactByConvID(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	react := ":cool:"
	channel := getOneOnOneChatChannel(t, "alice", "bob")

	lastMessageID := getMostRecentMessage(t, alice, channel).Id
	convID := getConvIDForChannel(t, alice, channel)

	// Send the react
	res, err := alice.ReactByConvID(convID, lastMessageID, react)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)
}

////////////////////////////////////////////////////////
// Manage channels /////////////////////////////////////
////////////////////////////////////////////////////////

func TestListChannels(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	teamChannel := getTeamChatChannel(t, "acme")
	channels, err := alice.ListChannels(teamChannel.Name)
	require.NoError(t, err)
	require.True(t, len(channels) > 0)
	channelFound := false
	for _, channel := range channels {
		if channel == teamChannel.TopicName {
			channelFound = true
			break
		}
	}
	require.True(t, channelFound)
}

func TestJoinAndLeaveChannel(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getTeamChatChannel(t, "acme")
	_, err := alice.LeaveChannel(channel.Name, channel.TopicName)
	require.NoError(t, err)
	_, err = alice.LeaveChannel(channel.Name, channel.TopicName)
	require.Error(t, err)
	_, err = alice.JoinChannel(channel.Name, channel.TopicName)
	require.NoError(t, err)
	_, err = alice.JoinChannel(channel.Name, channel.TopicName)
	// We don't get an error when trying to join an already joined oneOnOneChannel
	require.NoError(t, err)
}

////////////////////////////////////////////////////////
// Send lumens in chat /////////////////////////////////
////////////////////////////////////////////////////////

func TestInChatSend(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getOneOnOneChatChannel(t, "alice", "bob")
	text := "test InChatSend +1xlm " + randomString(t)

	// Send the message
	res, err := alice.InChatSend(channel, "%s", text)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(t, alice, channel)
	require.Equal(t, readMessage.Content.TypeName, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.Id, *res.Result.MessageID)
}

func TestInChatSendByConvID(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	text := "test InChatSendByConvID +1xlm " + randomString(t)
	channel := getOneOnOneChatChannel(t, "alice", "bob")
	convID := getConvIDForChannel(t, alice, channel)

	// Send the message
	res, err := alice.InChatSendByConvID(convID, "%s", text)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(t, alice, channel)
	require.Equal(t, readMessage.Content.TypeName, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.Id, *res.Result.MessageID)
}

func TestInChatSendByTlfName(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	text := "test InChatSendByTlfName +1xlm " + randomString(t)
	channel := getOneOnOneChatChannel(t, "alice", "bob")

	// Send the message
	res, err := alice.InChatSendByTlfName(channel.Name, "%s", text)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(t, alice, channel)
	require.Equal(t, readMessage.Content.TypeName, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.Id, *res.Result.MessageID)
}

// //////////////////////////////////////////////////////
// Misc commands ///////////////////////////////////////
// //////////////////////////////////////////////////////
func TestAdvertiseCommands(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)

	commands := []chat1.UserBotCommandInput{
		{
			Name:        "help",
			Description: "get help for a command",
			Usage:       "!help [cmdname]",
			ExtendedDescription: &chat1.UserBotExtendedDescription{
				Title:       "help command",
				DesktopBody: "help",
				MobileBody:  "help",
			},
		},
	}

	expectedOutput := []chat1.UserBotCommandOutput{
		{
			Username:    alice.GetUsername(),
			Name:        "help",
			Description: "get help for a command",
			Usage:       "!help [cmdname]",
			ExtendedDescription: &chat1.UserBotExtendedDescription{
				Title:       "help command",
				DesktopBody: "help",
				MobileBody:  "help",
			},
		},
	}

	_, err := alice.AdvertiseCommands(Advertisement{
		Alias: "botua",
		Advertisements: []chat1.AdvertiseCommandAPIParam{
			{
				Typ:      "public",
				Commands: commands,
			},
		},
	})
	require.NoError(t, err)

	teamChannel := getTeamChatChannel(t, "acme")
	res, err := alice.ListCommands(teamChannel)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, res)

	err = alice.ClearCommands(nil)
	require.NoError(t, err)

	res, err = alice.ListCommands(teamChannel)
	require.NoError(t, err)
	require.Zero(t, len(res))
}
