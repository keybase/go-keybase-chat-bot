package kbchat

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type keybaseTestConfig struct {
	Keybase string
	Bots    map[string]*OneshotOptions
	Teams   map[string]chat1.ChatChannel
}

func readAndParseTestConfig(t *testing.T) (config keybaseTestConfig) {
	data, err := ioutil.ReadFile("test_config.yaml")
	require.NoError(t, err)

	err = yaml.Unmarshal(data, &config)
	require.NoError(t, err)

	return config
}

func testBotSetup(t *testing.T, botName string) (bot *API, dir string) {
	config := readAndParseTestConfig(t)
	kbLocation := whichKeybase(t)
	if len(config.Keybase) > 0 {
		kbLocation = config.Keybase
	}
	dir = randomTempDir(t)
	kbTmpLocation := prepWorkingDir(t, dir, kbLocation)
	bot, err := Start(RunOptions{KeybaseLocation: kbTmpLocation, HomeDir: dir, Oneshot: config.Bots[botName], StartService: true})
	require.NoError(t, err)
	return bot, dir
}

func getOneOnOneChatChannel(t *testing.T, botName, oneOnOnePartner string) chat1.ChatChannel {
	config := readAndParseTestConfig(t)
	oneOnOneChannel := chat1.ChatChannel{
		Name: fmt.Sprintf("%s,%s", config.Bots[botName].Username, config.Bots[oneOnOnePartner].Username),
	}
	return oneOnOneChannel
}

func getTeamChatChannel(t *testing.T, teamName string) chat1.ChatChannel {
	config := readAndParseTestConfig(t)
	teamChannel := chat1.ChatChannel{
		Name:        config.Teams[teamName].Name,
		Public:      false,
		MembersType: "team",
		TopicName:   config.Teams[teamName].TopicName,
		TopicType:   "chat",
	}
	return teamChannel
}

func testBotTeardown(t *testing.T, bot *API, dir string) {
	err := bot.Shutdown()
	require.NoError(t, err)
	err = os.RemoveAll(dir)
	require.NoError(t, err)
}

func getMostRecentMessage(t *testing.T, bot *API, channel chat1.ChatChannel) chat1.MsgSummary {
	messages, err := bot.GetTextMessages(channel, false)
	require.NoError(t, err)
	return messages[0]
}

func getConvIDForChannel(t *testing.T, bot *API, channel chat1.ChatChannel) string {
	messages, err := bot.GetTextMessages(channel, false)
	require.NoError(t, err)
	convID := messages[0].ConvID
	return convID
}

func TestGetUsername(t *testing.T) {
	config := readAndParseTestConfig(t)
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	require.Equal(t, alice.GetUsername(), config.Bots["alice"].Username)
}

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
	res, err := alice.SendMessage(channel, text)
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
	res, err := alice.SendMessageByConvID(convID, text)
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
	res, err := alice.SendMessageByTlfName(channel.Name, text)
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
	res, err := alice.SendMessageByTeamName(channel.Name, &channel.TopicName, text)
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
	err := ioutil.WriteFile(location, data, 0644)
	require.NoError(t, err)

	// Send the message
	title := "test SendAttachmentByTeam " + randomString(t)
	res, err := alice.SendAttachmentByTeam(channel.Name, &channel.TopicName, location, title)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)
}

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

	err = alice.ClearCommands()
	require.NoError(t, err)

	res, err = alice.ListCommands(teamChannel)
	require.NoError(t, err)
	require.Zero(t, len(res))
}

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

func TestListenForNewTextMessages(t *testing.T) {
	alice, aliceDir := testBotSetup(t, "alice")
	bob, bobDir := testBotSetup(t, "bob")
	defer testBotTeardown(t, alice, aliceDir)
	defer testBotTeardown(t, bob, bobDir)
	channel := getOneOnOneChatChannel(t, "alice", "bob")

	sub, err := alice.ListenForNewTextMessages()
	require.NoError(t, err)

	done := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			message := strconv.Itoa(i)
			_, err := bob.SendMessage(channel, message)
			require.NoError(t, err)
		}
		done <- true
	}()

	receivedMessages := map[string]bool{
		"0": false,
		"1": false,
		"2": false,
		"3": false,
		"4": false,
	}

	for i := 0; i < 5; i++ {
		msg, err := sub.Read()
		require.NoError(t, err)
		require.Equal(t, msg.Message.Content.TypeName, "text")
		require.Equal(t, msg.Message.Sender.Username, bob.GetUsername())
		receivedMessages[msg.Message.Content.Text.Body] = true
	}

	for _, value := range receivedMessages {
		require.True(t, value)
	}

	<-done
}

func TestInChatSend(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getOneOnOneChatChannel(t, "alice", "bob")
	text := "test InChatSend +1xlm " + randomString(t)

	// Send the message
	res, err := alice.InChatSend(channel, text)
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
	res, err := alice.InChatSendByConvID(convID, text)
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
	res, err := alice.InChatSendByTlfName(channel.Name, text)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.True(t, *res.Result.MessageID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(t, alice, channel)
	require.Equal(t, readMessage.Content.TypeName, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.Id, *res.Result.MessageID)
}
