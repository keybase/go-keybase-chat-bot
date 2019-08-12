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

type botConfig struct {
	Bots  map[string]*OneshotOptions
	Teams map[string]Channel
}

func readAndParseConfig() botConfig {
	var config botConfig
	data, err := ioutil.ReadFile("test_config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func randomString() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func randomTempDir() string {
	dir := path.Join(os.TempDir(), "keybase_bot_"+randomString())
	return dir
}

func whichKeybase() (string, error) {
	cmd := exec.Command("which", "keybase")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	location := strings.TrimSpace(string(out))
	return location, nil
}

func copyFile(source, dest string) error {
	sourceData, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, sourceData, 0777)
	if err != nil {
		return err
	}

	return nil
}

// Creates the working directory and copies over the keybase binary in PATH.
// We do this to avoid any version mismatch issues.
func prepWorkingDir(workingDir string) (string, error) {
	kbLocation, err := whichKeybase()
	if err != nil {
		return "", err
	}

	err = os.Mkdir(workingDir, 0777)
	if err != nil {
		return "", err
	}
	kbDestination := path.Join(workingDir, "keybase")

	err = copyFile(kbLocation, kbDestination)
	if err != nil {
		return "", err
	}

	return kbDestination, nil
}

func deleteWorkingDir(workingDir string) error {
	return os.RemoveAll(workingDir)
}

type testSetupOptions struct {
	OneOnOnePartner string
	TeamName        string
}

func testSetup(botName string, options *testSetupOptions) (bot *API, config botConfig, dir string, oneOnOneChannel Channel, teamChannel Channel) {
	var oneOnOnePartner string
	var teamName string
	if options == nil {
		oneOnOnePartner = "charlie1"
		teamName = "acme"
	} else {
		oneOnOnePartner = options.OneOnOnePartner
		teamName = options.TeamName
	}

	config = readAndParseConfig()
	dir = randomTempDir()
	kbLocation, err := prepWorkingDir(dir)
	if err != nil {
		panic(err)
	}
	bot, err = Start(RunOptions{KeybaseLocation: kbLocation, HomeDir: dir, Oneshot: config.Bots[botName], StartService: true})
	if err != nil {
		panic(err)
	}

	oneOnOneChannel = Channel{
		Name: fmt.Sprintf("%s,%s", config.Bots[botName].Username, oneOnOnePartner),
	}
	teamChannel = Channel{
		Name:        config.Teams[teamName].Name,
		Public:      false,
		MembersType: "team",
		TopicName:   config.Teams[teamName].TopicName,
		TopicType:   "chat",
	}

	return bot, config, dir, oneOnOneChannel, teamChannel
}

func testTeardown(bot *API, dir string) {
	err := bot.Shutdown()
	if err != nil {
		panic(err)
	}
	err = deleteWorkingDir(dir)
	if err != nil {
		panic(err)
	}
}

func TestGetUsername(t *testing.T) {
	alice, config, dir, _, _ := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	require.Equal(t, alice.GetUsername(), config.Bots["alice1"].Username)
}

func TestGetConversations(t *testing.T) {
	alice, _, dir, _, _ := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	conversations, err := alice.GetConversations(false)
	require.NoError(t, err)
	require.Greater(t, len(conversations), 0)
}

func TestGetTextMessages(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	messages, err := alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	require.Greater(t, len(messages), 0)
}

func TestSendMessage(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	text := "test SendMessage " + randomString()

	// Send the message
	res, err := alice.SendMessage(oneOnOneChannel, text)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// Read it to confirm it sent
	messages, err := alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	sentMessage := messages[0]
	require.Equal(t, sentMessage.Content.Type, "text")
	require.Equal(t, sentMessage.Content.Text.Body, text)
	require.Equal(t, sentMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByConvID(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	text := "test SendMessageByConvID " + randomString()

	// Retrieve conversation ID
	messages, err := alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	convID := messages[0].ConversationID

	// Send the message
	res, err := alice.SendMessageByConvID(convID, text)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// Read it to confirm it sent
	messages, err = alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	sentMessage := messages[0]
	require.Equal(t, sentMessage.Content.Type, "text")
	require.Equal(t, sentMessage.Content.Text.Body, text)
	require.Equal(t, sentMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByTlfName(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	text := "test SendMessageByTlfName " + randomString()

	// Send the message
	res, err := alice.SendMessageByTlfName(oneOnOneChannel.Name, text)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// Read it to confirm it sent
	messages, err := alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	sentMessage := messages[0]
	require.Equal(t, sentMessage.Content.Type, "text")
	require.Equal(t, sentMessage.Content.Text.Body, text)
	require.Equal(t, sentMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByTeamName(t *testing.T) {
	alice, _, dir, _, teamChannel := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	text := "test SendMessageByTeamName " + randomString()

	// Send the message
	res, err := alice.SendMessageByTeamName(teamChannel.Name, text, &teamChannel.TopicName)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// Read it to confirm it sent
	messages, err := alice.GetTextMessages(teamChannel, false)
	require.NoError(t, err)
	sentMessage := messages[0]
	require.Equal(t, sentMessage.Content.Type, "text")
	require.Equal(t, sentMessage.Content.Text.Body, text)
	require.Equal(t, sentMessage.MsgID, res.Result.MsgID)
}

func TestSendAttachmentByTeam(t *testing.T) {
	alice, _, dir, _, teamChannel := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	// Create a test file
	fileName := "kb-attachment.txt"
	location := path.Join(os.TempDir(), fileName)
	data := []byte("My super cool attachment")
	err := ioutil.WriteFile(location, data, 0644)
	require.NoError(t, err)

	// Send the message
	title := "test SendAttachmentByTeam " + randomString()
	res, err := alice.SendAttachmentByTeam(teamChannel.Name, location, title, &teamChannel.TopicName)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// Types don't support attachments yet, so we can't read it
	// Read it to confirm it sent
	// messages, err := alice.GetTextMessages(teamChannel, false)
	// require.NoError(t, err)
	// sentMessage := messages[0]
	// require.Equal(t, sentMessage.Content.Type, "attachment")
	// require.Equal(t, sentMessage.Content.Attachment.Object.Title, title)
	// require.Equal(t, sentMessage.MsgID, res.Result.MsgID)
}

func TestReactByChannel(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	react := ":cool:"
	// Get last message, we'll react to it
	messages, err := alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	lastMessageID := messages[0].MsgID

	// Send the react
	res, err := alice.ReactByChannel(oneOnOneChannel, lastMessageID, react)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)
}

func TestReactByConvID(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	react := ":cool:"

	// Get last message, we'll react to it
	messages, err := alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	lastMessageID := messages[0].MsgID

	// Retrieve conversation ID
	messages, err = alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	convID := messages[0].ConversationID

	// Send the react
	res, err := alice.ReactByConvID(convID, lastMessageID, react)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)
}

func TestAdvertiseCommands(t *testing.T) {}

func TestListChannels(t *testing.T) {
	alice, _, dir, _, teamChannel := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	channels, err := alice.ListChannels(teamChannel.Name)
	require.NoError(t, err)
	require.Greater(t, len(channels), 0)
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
	alice, _, dir, _, teamChannel := testSetup("alice1", nil)
	defer testTeardown(alice, dir)
	_, err := alice.LeaveChannel(teamChannel.Name, teamChannel.TopicName)
	require.NoError(t, err)
	_, err = alice.LeaveChannel(teamChannel.Name, teamChannel.TopicName)
	require.Error(t, err)
	_, err = alice.JoinChannel(teamChannel.Name, teamChannel.TopicName)
	require.NoError(t, err)
	_, err = alice.JoinChannel(teamChannel.Name, teamChannel.TopicName)
	// We don't get an error when trying to join an already joined oneOnOneChannel
	require.NoError(t, err)
}

func TestListenForNewTextMessages(t *testing.T) {
	alice, _, aliceDir, oneOnOneChannel, _ := testSetup("alice1", nil)
	bob, _, bobDir, _, _ := testSetup("bob1", nil)

	sub, err := alice.ListenForNewTextMessages()
	require.NoError(t, err)

	go func() {
		defer testTeardown(alice, aliceDir)
		defer testTeardown(bob, bobDir)

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
	}()

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		message := strconv.Itoa(i)
		_, err := bob.SendMessage(oneOnOneChannel, message)
		require.NoError(t, err)
	}
}
