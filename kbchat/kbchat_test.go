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

type team struct {
	Teamname string
	Channel  string
}

type botConfig struct {
	Bots  map[string]*OneshotOptions
	Teams map[string]team
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

func testSetup() (alice *API, config botConfig, dir string, oneOnOneChannel Channel, teamChannel Channel) {
	var err error
	config = readAndParseConfig()
	dir = randomTempDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in generating directory: %v\n", err)
	}
	kbLocation, err := prepWorkingDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in preparing working directory: %v\n", err)
	}
	alice, err = Start(RunOptions{KeybaseLocation: kbLocation, HomeDir: dir, Oneshot: config.Bots["alice1"], StartService: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in starting service: %v\n", err)
	}

	oneOnOneChannel = Channel{
		Name: fmt.Sprintf("%s,%s", config.Bots["alice1"].Username, config.Bots["charlie1"].Username),
	}
	teamChannel = Channel{
		Name:        config.Teams["acme"].Teamname,
		Public:      false,
		MembersType: "team",
		TopicName:   config.Teams["acme"].Channel,
		TopicType:   "chat",
	}

	return alice, config, dir, oneOnOneChannel, teamChannel
}

func testTeardown(alice *API, dir string) {
	err := alice.Shutdown()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on service shutdown: %v\n", err)
	}
	err = deleteWorkingDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error working directory teardown: %v\n", err)
	}
}

func TestGetUsername(t *testing.T) {
	alice, config, dir, _, _ := testSetup()
	require.Equal(t, alice.GetUsername(), config.Bots["alice1"].Username)
	testTeardown(alice, dir)
}

func TestGetConversations(t *testing.T) {
	alice, _, dir, _, _ := testSetup()
	conversations, err := alice.GetConversations(false)
	require.NoError(t, err)
	require.Greater(t, len(conversations), 0)
	testTeardown(alice, dir)
}

func TestGetTextMessages(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup()
	messages, err := alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	require.Greater(t, len(messages), 0)
	testTeardown(alice, dir)
}

func TestSendMessage(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup()
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
	testTeardown(alice, dir)
}

func TestSendMessageByConvID(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup()
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
	testTeardown(alice, dir)
}

func TestSendMessageByTlfName(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup()
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
	testTeardown(alice, dir)
}

func TestSendMessageByTeamName(t *testing.T) {
	alice, _, dir, _, teamChannel := testSetup()
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
	testTeardown(alice, dir)
}

func TestSendAttachmentByTeam(t *testing.T) {
	alice, _, dir, _, teamChannel := testSetup()
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
	testTeardown(alice, dir)
}

func TestReactByChannel(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup()
	react := ":cool:"
	// Get last message, we'll react to it
	messages, err := alice.GetTextMessages(oneOnOneChannel, false)
	require.NoError(t, err)
	lastMessageID := messages[0].MsgID

	// Send the react
	res, err := alice.ReactByChannel(oneOnOneChannel, lastMessageID, react)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// No great way to confirm reaction yet
	testTeardown(alice, dir)
}

func TestReactByConvID(t *testing.T) {
	alice, _, dir, oneOnOneChannel, _ := testSetup()
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
	testTeardown(alice, dir)
}

func TestAdvertiseCommands(t *testing.T) {}

func TestListChannels(t *testing.T) {
	alice, _, dir, _, teamChannel := testSetup()
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
	testTeardown(alice, dir)
}

func TestJoinAndLeaveChannel(t *testing.T) {
	alice, _, dir, _, teamChannel := testSetup()
	_, err := alice.LeaveChannel(teamChannel.Name, teamChannel.TopicName)
	require.NoError(t, err)
	_, err = alice.LeaveChannel(teamChannel.Name, teamChannel.TopicName)
	require.Error(t, err)
	_, err = alice.JoinChannel(teamChannel.Name, teamChannel.TopicName)
	require.NoError(t, err)
	_, err = alice.JoinChannel(teamChannel.Name, teamChannel.TopicName)
	// We don't get an error when trying to join an already joined oneOnOneChannel
	require.NoError(t, err)
	testTeardown(alice, dir)
}

func TestListenForNewTextMessages(t *testing.T) {
	alice, config, dir, oneOnOneChannel, _ := testSetup()
	bobDir := randomTempDir()
	kbLocation, err := prepWorkingDir(bobDir)
	require.NoError(t, err)
	bob, err := Start(RunOptions{KeybaseLocation: kbLocation, HomeDir: bobDir, Oneshot: config.Bots["bob1"], StartService: true})
	require.NoError(t, err)

	sub, err := alice.ListenForNewTextMessages()
	require.NoError(t, err)

	go func() {
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

		testTeardown(alice, dir)
	}()

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		message := strconv.Itoa(i)
		_, err := bob.SendMessage(oneOnOneChannel, message)
		require.NoError(t, err)
	}

}
