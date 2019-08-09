package kbchat

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
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

type bots struct {
	Alice1   *OneshotOptions
	Alice2   *OneshotOptions
	Bob1     *OneshotOptions
	Charlie1 *OneshotOptions
}

type team struct {
	Teamname string
	Channel  string
}

type teams struct {
	Acme             team
	AlicesPlayground team
}

type botConfig struct {
	Bots  bots
	Teams teams
}

func readAndParseConfig() (botConfig, error) {
	var config botConfig
	data, err := ioutil.ReadFile("test_config.yaml")
	if err != nil {
		return botConfig{}, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return botConfig{}, err
	}

	return config, nil
}

func randomTempDir() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	dir := path.Join(os.TempDir(), "keybase_bot_"+hex.EncodeToString(bytes))
	return dir, nil
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

var alice *API
var config botConfig
var channel Channel
var teamChannel Channel

func TestMain(m *testing.M) {
	var err error
	config, err = readAndParseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in reading config: %v\n", err)
	}
	dir, err := randomTempDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in generating directory: %v\n", err)
	}
	kbLocation, err := prepWorkingDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in preparing working directory: %v\n", err)
	}
	alice, err = Start(RunOptions{KeybaseLocation: kbLocation, HomeDir: dir, Oneshot: config.Bots.Alice1, StartService: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in starting service: %v\n", err)
	}

	channel = Channel{
		Name: fmt.Sprintf("%s,%s", config.Bots.Alice1.Username, config.Bots.Charlie1.Username),
	}
	teamChannel = Channel{
		Name:        config.Teams.Acme.Teamname,
		Public:      false,
		MembersType: "team",
		TopicName:   config.Teams.Acme.Channel,
		TopicType:   "chat",
	}

	flag.Parse()
	code := m.Run()

	err = alice.Shutdown()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on service shutdown: %v\n", err)
	}
	err = deleteWorkingDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error working directory teardown: %v\n", err)
	}
	os.Exit(code)
}

func TestGetUsername(t *testing.T) {
	require.Equal(t, alice.GetUsername(), config.Bots.Alice1.Username)
}

func TestGetConversations(t *testing.T) {
	conversations, err := alice.GetConversations(false)
	require.NoError(t, err)
	require.Greater(t, len(conversations), 0)
}

func TestGetTextMessages(t *testing.T) {
	messages, err := alice.GetTextMessages(channel, false)
	require.NoError(t, err)
	require.Greater(t, len(messages), 0)
}

func TestSendMessage(t *testing.T) {
	text := "test SendMessage"

	// Send the message
	res, err := alice.SendMessage(channel, text)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// Read it to confirm it sent
	messages, err := alice.GetTextMessages(channel, false)
	require.NoError(t, err)
	sentMessage := messages[0]
	require.Equal(t, sentMessage.Content.Type, "text")
	require.Equal(t, sentMessage.Content.Text.Body, text)
	require.Equal(t, sentMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByConvID(t *testing.T) {
	text := "test SendMessageByConvID"

	// Retrieve conversation ID
	messages, err := alice.GetTextMessages(channel, false)
	convID := messages[0].ConversationID

	// Send the message
	res, err := alice.SendMessageByConvID(convID, text)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// Read it to confirm it sent
	messages, err = alice.GetTextMessages(channel, false)
	require.NoError(t, err)
	sentMessage := messages[0]
	require.Equal(t, sentMessage.Content.Type, "text")
	require.Equal(t, sentMessage.Content.Text.Body, text)
	require.Equal(t, sentMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByTlfName(t *testing.T) {
	text := "test SendMessageByTlfName"

	// Send the message
	res, err := alice.SendMessageByTlfName(channel.Name, text)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// Read it to confirm it sent
	messages, err := alice.GetTextMessages(channel, false)
	require.NoError(t, err)
	sentMessage := messages[0]
	require.Equal(t, sentMessage.Content.Type, "text")
	require.Equal(t, sentMessage.Content.Text.Body, text)
	require.Equal(t, sentMessage.MsgID, res.Result.MsgID)
}

func TestSendMessageByTeamName(t *testing.T) {
	text := "test SendMessageByTeamName"

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
	// Create a test file
	fileName := "kb-attachment.txt"
	location := path.Join(os.TempDir(), fileName)
	data := []byte("My super cool attachment")
	err := ioutil.WriteFile(location, data, 0644)
	require.NoError(t, err)

	// Send the message
	title := "test SendAttachmentByTeam"
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
	react := ":cool:"
	// Get last message, we'll react to it
	messages, err := alice.GetTextMessages(channel, false)
	require.NoError(t, err)
	lastMessageID := messages[0].MsgID

	// Send the react
	res, err := alice.ReactByChannel(channel, lastMessageID, react)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)

	// No great way to confirm reaction yet
}

func TestReactByConvID(t *testing.T) {
	react := ":cool:"

	// Get last message, we'll react to it
	messages, err := alice.GetTextMessages(channel, false)
	require.NoError(t, err)
	lastMessageID := messages[0].MsgID

	// Retrieve conversation ID
	messages, err = alice.GetTextMessages(channel, false)
	require.NoError(t, err)
	convID := messages[0].ConversationID

	// Send the react
	res, err := alice.ReactByConvID(convID, lastMessageID, react)
	require.NoError(t, err)
	require.Greater(t, res.Result.MsgID, 0)
}

func TestAdvertiseCommands(t *testing.T) {}

func TestListChannels(t *testing.T) {
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
	_, err := alice.LeaveChannel(teamChannel.Name, teamChannel.TopicName)
	require.NoError(t, err)
	_, err = alice.LeaveChannel(teamChannel.Name, teamChannel.TopicName)
	require.Error(t, err)
	_, err = alice.JoinChannel(teamChannel.Name, teamChannel.TopicName)
	require.NoError(t, err)
	_, err = alice.JoinChannel(teamChannel.Name, teamChannel.TopicName)
	// We don't get an error when trying to join an already joined channel
	require.NoError(t, err)
}

func TestListenForNewTextMessages(t *testing.T) {
	dir, err := randomTempDir()
	require.NoError(t, err)
	kbLocation, err := prepWorkingDir(dir)
	require.NoError(t, err)
	bob, err := Start(RunOptions{KeybaseLocation: kbLocation, HomeDir: dir, Oneshot: config.Bots.Bob1, StartService: true})
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
	}()

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		message := strconv.Itoa(i)
		_, err := bob.SendMessage(channel, message)
		require.NoError(t, err)
	}
}
