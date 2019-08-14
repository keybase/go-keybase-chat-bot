package kbchat

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type keybaseTestConfig struct {
	Bots  map[string]*OneshotOptions
	Teams map[string]Channel
}

func readAndParseTestConfig(t *testing.T) (config keybaseTestConfig) {
	data, err := ioutil.ReadFile("test_config.yaml")
	require.NoError(t, err)

	err = yaml.Unmarshal(data, &config)
	require.NoError(t, err)

	return config
}

func randomString(t *testing.T) string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	require.NoError(t, err)
	return hex.EncodeToString(bytes)
}

func randomTempDir(t *testing.T) string {
	return path.Join(os.TempDir(), "keybase_bot_"+randomString(t))
}

func whichKeybase(t *testing.T) string {
	cmd := exec.Command("which", "keybase")
	out, err := cmd.Output()
	require.NoError(t, err)
	location := strings.TrimSpace(string(out))
	return location
}

func copyFile(t *testing.T, source, dest string) {
	sourceData, err := ioutil.ReadFile(source)
	require.NoError(t, err)
	err = ioutil.WriteFile(dest, sourceData, 0777)
	require.NoError(t, err)
}

// Creates the working directory and copies over the keybase binary in PATH.
// We do this to avoid any version mismatch issues.
func prepWorkingDir(t *testing.T, workingDir string) string {
	kbLocation := whichKeybase(t)

	err := os.Mkdir(workingDir, 0777)
	require.NoError(t, err)
	kbDestination := path.Join(workingDir, "keybase")

	copyFile(t, kbLocation, kbDestination)

	return kbDestination
}

func testBotSetup(t *testing.T, botName string) (bot *API, dir string) {
	config := readAndParseTestConfig(t)
	dir = randomTempDir(t)
	kbLocation := prepWorkingDir(t, dir)
	bot, err := Start(RunOptions{KeybaseLocation: kbLocation, HomeDir: dir, Oneshot: config.Bots[botName], StartService: true})
	require.NoError(t, err)
	return bot, dir
}

func getOneOnOneChatChannel(t *testing.T, botName, oneOnOnePartner string) Channel {
	config := readAndParseTestConfig(t)
	oneOnOneChannel := Channel{
		Name: fmt.Sprintf("%s,%s", config.Bots[botName].Username, config.Bots[oneOnOnePartner].Username),
	}
	return oneOnOneChannel
}

func getTeamChatChannel(t *testing.T, teamName string) Channel {
	config := readAndParseTestConfig(t)
	teamChannel := Channel{
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

func getMostRecentMessage(bot *API, channel Channel) Message {
	messages, err := bot.GetTextMessages(channel, false)
	if err != nil {
		panic(err)
	}
	return messages[0]
}

func getConvIDForChannel(bot *API, channel Channel) string {
	messages, err := bot.GetTextMessages(channel, false)
	if err != nil {
		panic(err)
	}
	convID := messages[0].ConversationID
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
	require.True(t, res.Result.MsgID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(alice, channel)
	require.Equal(t, readMessage.Content.Type, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByConvID(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	text := "test SendMessageByConvID " + randomString(t)
	channel := getOneOnOneChatChannel(t, "alice", "bob")
	convID := getConvIDForChannel(alice, channel)

	// Send the message
	res, err := alice.SendMessageByConvID(convID, text)
	require.NoError(t, err)
	require.True(t, res.Result.MsgID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(alice, channel)
	require.Equal(t, readMessage.Content.Type, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByTlfName(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	text := "test SendMessageByTlfName " + randomString(t)
	channel := getOneOnOneChatChannel(t, "alice", "bob")

	// Send the message
	res, err := alice.SendMessageByTlfName(channel.Name, text)
	require.NoError(t, err)
	require.True(t, res.Result.MsgID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(alice, channel)
	require.Equal(t, readMessage.Content.Type, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByTeamName(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	text := "test SendMessageByTeamName " + randomString(t)
	channel := getTeamChatChannel(t, "acme")

	// Send the message
	res, err := alice.SendMessageByTeamName(channel.Name, text, &channel.TopicName)
	require.NoError(t, err)
	require.True(t, res.Result.MsgID > 0)

	// Read it to confirm it sent
	readMessage := getMostRecentMessage(alice, channel)
	require.Equal(t, readMessage.Content.Type, "text")
	require.Equal(t, readMessage.Content.Text.Body, text)
	require.Equal(t, readMessage.MsgID, res.Result.MsgID)
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
	res, err := alice.SendAttachmentByTeam(channel.Name, location, title, &channel.TopicName)
	require.NoError(t, err)
	require.True(t, res.Result.MsgID > 0)
}

func TestReactByChannel(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	channel := getOneOnOneChatChannel(t, "alice", "bob")

	react := ":cool:"
	lastMessageID := getMostRecentMessage(alice, channel).MsgID

	res, err := alice.ReactByChannel(channel, lastMessageID, react)
	require.NoError(t, err)
	require.True(t, res.Result.MsgID > 0)
}

func TestReactByConvID(t *testing.T) {
	alice, dir := testBotSetup(t, "alice")
	defer testBotTeardown(t, alice, dir)
	react := ":cool:"
	channel := getOneOnOneChatChannel(t, "alice", "bob")

	lastMessageID := getMostRecentMessage(alice, channel).MsgID
	convID := getConvIDForChannel(alice, channel)

	// Send the react
	res, err := alice.ReactByConvID(convID, lastMessageID, react)
	require.NoError(t, err)
	require.True(t, res.Result.MsgID > 0)
}

func TestAdvertiseCommands(t *testing.T) {}

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
		require.Equal(t, msg.Message.Content.Type, "text")
		require.Equal(t, msg.Message.Sender.Username, bob.GetUsername())
		receivedMessages[msg.Message.Content.Text.Body] = true
	}

	for _, value := range receivedMessages {
		require.True(t, value)
	}

	<-done
}
