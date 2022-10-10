package kbchat

import (
	"fmt"
	"os"
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
	data, err := os.ReadFile("test_config.yaml")
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

func getConvIDForChannel(t *testing.T, bot *API, channel chat1.ChatChannel) chat1.ConvIDStr {
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
