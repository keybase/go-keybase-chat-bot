package kbchat

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/stellar1"
)

// API is the main object used for communicating with the Keybase JSON API
type API struct {
	sync.Mutex
	apiInput  io.Writer
	apiOutput *bufio.Reader
	apiCmd    *exec.Cmd
	username  string
	runOpts   RunOptions
}

func getUsername(runOpts RunOptions) (username string, err error) {
	p := runOpts.Command("status")
	output, err := p.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = p.Start(); err != nil {
		return "", err
	}

	doneCh := make(chan error)
	go func() {
		scanner := bufio.NewScanner(output)
		if !scanner.Scan() {
			doneCh <- errors.New("unable to find Keybase username")
			return
		}
		toks := strings.Fields(scanner.Text())
		if len(toks) != 2 {
			doneCh <- errors.New("invalid Keybase username output")
			return
		}
		username = toks[1]
		doneCh <- nil
	}()

	select {
	case err = <-doneCh:
		if err != nil {
			return "", err
		}
	case <-time.After(5 * time.Second):
		return "", errors.New("unable to run Keybase command")
	}

	return username, nil
}

type OneshotOptions struct {
	Username string
	PaperKey string
}

type RunOptions struct {
	KeybaseLocation string
	HomeDir         string
	Oneshot         *OneshotOptions
	StartService    bool
}

func (r RunOptions) Location() string {
	if r.KeybaseLocation == "" {
		return "keybase"
	}
	return r.KeybaseLocation
}

func (r RunOptions) Command(args ...string) *exec.Cmd {
	var cmd []string
	if r.HomeDir != "" {
		cmd = append(cmd, "--home", r.HomeDir)
	}
	cmd = append(cmd, args...)
	return exec.Command(r.Location(), cmd...)
}

// Start fires up the Keybase JSON API in stdin/stdout mode
func Start(runOpts RunOptions) (*API, error) {
	api := &API{
		runOpts: runOpts,
	}
	if err := api.startPipes(); err != nil {
		return nil, err
	}
	return api, nil
}

func (a *API) auth() (string, error) {
	username, err := getUsername(a.runOpts)
	if err == nil {
		return username, nil
	}
	if a.runOpts.Oneshot == nil {
		return "", err
	}
	username = ""
	// If a paper key is specified, then login with oneshot mode (logout first)
	if a.runOpts.Oneshot != nil {
		if username == a.runOpts.Oneshot.Username {
			// just get out if we are on the desired user already
			return username, nil
		}
		if err := a.runOpts.Command("logout", "-f").Run(); err != nil {
			return "", err
		}
		if err := a.runOpts.Command("oneshot", "--username", a.runOpts.Oneshot.Username, "--paperkey",
			a.runOpts.Oneshot.PaperKey).Run(); err != nil {
			return "", err
		}
		username = a.runOpts.Oneshot.Username
		return username, nil
	}
	return "", errors.New("unable to auth")
}

func (a *API) startPipes() (err error) {
	a.Lock()
	defer a.Unlock()
	if a.apiCmd != nil {
		if err := a.apiCmd.Process.Kill(); err != nil {
			return err
		}
	}
	a.apiCmd = nil

	if a.runOpts.StartService {
		if err := a.runOpts.Command("service").Start(); err != nil {
			return err
		}
	}

	if a.username, err = a.auth(); err != nil {
		return err
	}
	a.apiCmd = a.runOpts.Command("chat", "api")
	if a.apiInput, err = a.apiCmd.StdinPipe(); err != nil {
		return err
	}
	output, err := a.apiCmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := a.apiCmd.Start(); err != nil {
		return err
	}
	a.apiOutput = bufio.NewReader(output)
	return nil
}

var errAPIDisconnected = errors.New("chat API disconnected")

func (a *API) getAPIPipesLocked() (io.Writer, *bufio.Reader, error) {
	// this should only be called inside a lock
	if a.apiCmd == nil {
		return nil, nil, errAPIDisconnected
	}
	return a.apiInput, a.apiOutput, nil
}

// GetConversations reads all conversations from the current user's inbox.
func (a *API) GetConversations(unreadOnly bool) ([]chat1.ConvSummary, error) {
	apiInput := fmt.Sprintf(`{"method":"list", "params": { "options": { "unread_only": %v}}}`, unreadOnly)
	output, err := a.doFetch(apiInput)
	if err != nil {
		return nil, err
	}

	var inbox Inbox
	if err := json.Unmarshal(output, &inbox); err != nil {
		return nil, err
	} else if inbox.Error != nil {
		return nil, errors.New(inbox.Error.Message)
	}
	return inbox.Result.Convs, nil
}

// GetTextMessages fetches all text messages from a given channel. Optionally can filter
// ont unread status.
func (a *API) GetTextMessages(channel chat1.ChatChannel, unreadOnly bool) ([]chat1.MsgSummary, error) {
	channelBytes, err := json.Marshal(channel)
	if err != nil {
		return nil, err
	}
	apiInput := fmt.Sprintf(`{"method": "read", "params": {"options": {"channel": %s}}}`, string(channelBytes))
	output, err := a.doFetch(apiInput)
	if err != nil {
		return nil, err
	}

	var thread Thread

	if err := json.Unmarshal(output, &thread); err != nil {
		return nil, fmt.Errorf("unable to decode thread: %v", err)
	} else if thread.Error != nil {
		return nil, errors.New(thread.Error.Message)
	}

	var res []chat1.MsgSummary
	for _, msg := range thread.Result.Messages {
		if msg.Msg.Content.TypeName == "text" {
			res = append(res, *msg.Msg)
		}
	}

	return res, nil
}

type sendMessageBody struct {
	Body string
}

type sendMessageOptions struct {
	Channel          chat1.ChatChannel `json:"channel,omitempty"`
	ConversationID   string            `json:"conversation_id,omitempty"`
	Message          sendMessageBody   `json:",omitempty"`
	Filename         string            `json:"filename,omitempty"`
	Title            string            `json:"title,omitempty"`
	MsgID            chat1.MessageID   `json:"message_id,omitempty"`
	ConfirmLumenSend bool              `json:"confirm_lumen_send"`
}

type sendMessageParams struct {
	Options sendMessageOptions
}

type sendMessageArg struct {
	Method string
	Params sendMessageParams
}

func (a *API) doSend(arg interface{}) (resp SendResponse, err error) {
	a.Lock()
	defer a.Unlock()

	bArg, err := json.Marshal(arg)
	if err != nil {
		return SendResponse{}, err
	}
	input, output, err := a.getAPIPipesLocked()
	if err != nil {
		return SendResponse{}, err
	}
	if _, err := io.WriteString(input, string(bArg)); err != nil {
		return SendResponse{}, err
	}
	responseRaw, err := output.ReadBytes('\n')
	if err != nil {
		return SendResponse{}, err
	}
	if err := json.Unmarshal(responseRaw, &resp); err != nil {
		return resp, fmt.Errorf("failed to decode API response: %s", err)
	} else if resp.Error != nil {
		return resp, errors.New(resp.Error.Message)
	}
	return resp, nil
}

func (a *API) doFetch(apiInput string) ([]byte, error) {
	a.Lock()
	defer a.Unlock()

	input, output, err := a.getAPIPipesLocked()
	if err != nil {
		return nil, err
	}
	if _, err := io.WriteString(input, apiInput); err != nil {
		return nil, err
	}
	byteOutput, err := output.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return byteOutput, nil
}

func newSendArg(options sendMessageOptions) sendMessageArg {
	return sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: options,
		},
	}
}

func (a *API) SendMessage(channel chat1.ChatChannel, body string) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		Channel: channel,
		Message: sendMessageBody{
			Body: body,
		},
	})
	return a.doSend(arg)
}

func (a *API) SendMessageByConvID(convID string, body string) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		ConversationID: convID,
		Message: sendMessageBody{
			Body: body,
		},
	})
	return a.doSend(arg)
}

// SendMessageByTlfName sends a message on the given TLF name
func (a *API) SendMessageByTlfName(tlfName string, body string) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		Channel: chat1.ChatChannel{
			Name: tlfName,
		},
		Message: sendMessageBody{
			Body: body,
		},
	})
	return a.doSend(arg)
}

func (a *API) SendMessageByTeamName(teamName string, body string, inChannel *string) (SendResponse, error) {
	channel := "general"
	if inChannel != nil {
		channel = *inChannel
	}
	arg := newSendArg(sendMessageOptions{
		Channel: chat1.ChatChannel{
			MembersType: "team",
			Name:        teamName,
			TopicName:   channel,
		},
		Message: sendMessageBody{
			Body: body,
		},
	})
	return a.doSend(arg)
}

func (a *API) SendAttachmentByTeam(teamName string, filename string, title string, inChannel *string) (SendResponse, error) {
	channel := "general"
	if inChannel != nil {
		channel = *inChannel
	}
	arg := sendMessageArg{
		Method: "attach",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				Channel: chat1.ChatChannel{
					MembersType: "team",
					Name:        teamName,
					TopicName:   channel,
				},
				Filename: filename,
				Title:    title,
			},
		},
	}
	return a.doSend(arg)
}

type reactionOptions struct {
	ConversationID string `json:"conversation_id"`
	Message        sendMessageBody
	MsgID          chat1.MessageID   `json:"message_id"`
	Channel        chat1.ChatChannel `json:"channel"`
}

type reactionParams struct {
	Options reactionOptions
}

type reactionArg struct {
	Method string
	Params reactionParams
}

func newReactionArg(options reactionOptions) reactionArg {
	return reactionArg{
		Method: "reaction",
		Params: reactionParams{Options: options},
	}
}

func (a *API) ReactByChannel(channel chat1.ChatChannel, msgID chat1.MessageID, reaction string) (SendResponse, error) {
	arg := newReactionArg(reactionOptions{
		Message: sendMessageBody{Body: reaction},
		MsgID:   msgID,
		Channel: channel,
	})
	return a.doSend(arg)
}

func (a *API) ReactByConvID(convID string, msgID chat1.MessageID, reaction string) (SendResponse, error) {
	arg := newReactionArg(reactionOptions{
		Message:        sendMessageBody{Body: reaction},
		MsgID:          msgID,
		ConversationID: convID,
	})
	return a.doSend(arg)
}

func (a *API) InChatSend(channel chat1.ChatChannel, body string) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		Channel: channel,
		Message: sendMessageBody{
			Body: body,
		},
		ConfirmLumenSend: true,
	})
	return a.doSend(arg)
}

func (a *API) InChatSendByConvID(convID string, body string) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		ConversationID: convID,
		Message: sendMessageBody{
			Body: body,
		},
		ConfirmLumenSend: true,
	})
	return a.doSend(arg)
}

func (a *API) InChatSendByTlfName(tlfName string, body string) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		Channel: chat1.ChatChannel{
			Name: tlfName,
		},
		Message: sendMessageBody{
			Body: body,
		},
		ConfirmLumenSend: true,
	})
	return a.doSend(arg)
}

type advertiseCmdsParams struct {
	Options Advertisement
}

type advertiseCmdsMsgArg struct {
	Method string
	Params advertiseCmdsParams
}

func newAdvertiseCmdsMsgArg(ad Advertisement) advertiseCmdsMsgArg {
	return advertiseCmdsMsgArg{
		Method: "advertisecommands",
		Params: advertiseCmdsParams{
			Options: ad,
		},
	}
}

func (a *API) AdvertiseCommands(ad Advertisement) (SendResponse, error) {
	return a.doSend(newAdvertiseCmdsMsgArg(ad))
}

func (a *API) ClearCommands() error {
	arg := struct {
		Method string
	}{
		Method: "clearcommands",
	}
	_, err := a.doSend(arg)
	return err
}

type listCmdsOptions struct {
	Channel        chat1.ChatChannel
	ConversationID string
}

type listCmdsParams struct {
	Options listCmdsOptions
}

type listCmdsArg struct {
	Method string
	Params listCmdsParams
}

func newListCmdsArg(options listCmdsOptions) listCmdsArg {
	return listCmdsArg{
		Method: "listcommands",
		Params: listCmdsParams{
			Options: options,
		},
	}
}

func (a *API) ListCommands(channel chat1.ChatChannel) ([]chat1.UserBotCommandOutput, error) {
	arg := newListCmdsArg(listCmdsOptions{
		Channel: channel,
	})
	return a.listCommands(arg)
}

func (a *API) ListCommandsByConvID(conversationID string) ([]chat1.UserBotCommandOutput, error) {
	arg := newListCmdsArg(listCmdsOptions{
		ConversationID: conversationID,
	})
	return a.listCommands(arg)
}

func (a *API) listCommands(arg listCmdsArg) ([]chat1.UserBotCommandOutput, error) {
	bArg, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}
	output, err := a.doFetch(string(bArg))
	if err != nil {
		return nil, err
	}
	var res ListCommandsResponse
	if err := json.Unmarshal(output, &res); err != nil {
		return nil, err
	} else if res.Error != nil {
		return nil, errors.New(res.Error.Message)
	}
	return res.Result.Commands, nil
}

// SubscriptionMessage contains a message and conversation object
type SubscriptionMessage struct {
	Message      chat1.MsgSummary
	Conversation chat1.ConvSummary
}

type SubscriptionWalletEvent struct {
	Payment stellar1.PaymentDetailsLocal
}

// NewSubscription has methods to control the background message fetcher loop
type NewSubscription struct {
	newMsgsCh   <-chan SubscriptionMessage
	newWalletCh <-chan SubscriptionWalletEvent
	errorCh     <-chan error
	shutdownCh  chan struct{}
}

// Read blocks until a new message arrives
func (m NewSubscription) Read() (SubscriptionMessage, error) {
	select {
	case msg := <-m.newMsgsCh:
		return msg, nil
	case err := <-m.errorCh:
		return SubscriptionMessage{}, err
	}
}

// Read blocks until a new message arrives
func (m NewSubscription) ReadWallet() (SubscriptionWalletEvent, error) {
	select {
	case msg := <-m.newWalletCh:
		return msg, nil
	case err := <-m.errorCh:
		return SubscriptionWalletEvent{}, err
	}
}

// Shutdown terminates the background process
func (m NewSubscription) Shutdown() {
	m.shutdownCh <- struct{}{}
}

type ListenOptions struct {
	Wallet bool
}

// ListenForNewTextMessages proxies to Listen without wallet events
func (a *API) ListenForNewTextMessages() (NewSubscription, error) {
	opts := ListenOptions{Wallet: false}
	return a.Listen(opts)
}

// Listen fires of a background loop and puts chat messages and wallet
// events into channels
func (a *API) Listen(opts ListenOptions) (NewSubscription, error) {
	newMsgCh := make(chan SubscriptionMessage, 100)
	newWalletCh := make(chan SubscriptionWalletEvent, 100)
	errorCh := make(chan error, 100)
	shutdownCh := make(chan struct{})
	done := make(chan struct{})

	sub := NewSubscription{
		newMsgsCh:   newMsgCh,
		newWalletCh: newWalletCh,
		shutdownCh:  shutdownCh,
		errorCh:     errorCh,
	}
	pause := 2 * time.Second
	readScanner := func(boutput *bufio.Scanner) {
		for {
			boutput.Scan()
			t := boutput.Text()
			var typeHolder TypeHolder
			if err := json.Unmarshal([]byte(t), &typeHolder); err != nil {
				errorCh <- err
				break
			}
			switch typeHolder.Type {
			case "chat":
				var notification chat1.MsgNotification
				if err := json.Unmarshal([]byte(t), &notification); err != nil {
					errorCh <- err
					break
				}
				subscriptionMessage := SubscriptionMessage{
					Message: *notification.Msg,
					Conversation: chat1.ConvSummary{
						Id:      notification.Msg.ConvID,
						Channel: notification.Msg.Channel,
					},
				}
				newMsgCh <- subscriptionMessage
			case "wallet":
				var holder PaymentHolder
				if err := json.Unmarshal([]byte(t), &holder); err != nil {
					errorCh <- err
					break
				}
				subscriptionPayment := SubscriptionWalletEvent(holder)
				newWalletCh <- subscriptionPayment
			default:
				continue
			}
		}
		done <- struct{}{}
	}

	attempts := 0
	maxAttempts := 1800
	go func() {
		for {
			if attempts >= maxAttempts {
				panic("Listen: failed to auth, giving up")
			}
			attempts++
			if _, err := a.auth(); err != nil {
				log.Printf("Listen: failed to auth: %s", err)
				time.Sleep(pause)
				continue
			}
			cmdElements := []string{"chat", "api-listen"}
			if opts.Wallet {
				cmdElements = append(cmdElements, "--wallet")
			}
			p := a.runOpts.Command(cmdElements...)
			output, err := p.StdoutPipe()
			if err != nil {
				log.Printf("Listen: failed to listen: %s", err)
				time.Sleep(pause)
				continue
			}
			boutput := bufio.NewScanner(output)
			if err := p.Start(); err != nil {
				log.Printf("Listen: failed to make listen scanner: %s", err)
				time.Sleep(pause)
				continue
			}
			attempts = 0
			go readScanner(boutput)
			<-done
			if err := p.Wait(); err != nil {
				log.Printf("Listen: failed to Wait for command: %s", err)
			}
			time.Sleep(pause)
		}
	}()
	return sub, nil
}

func (a *API) GetUsername() string {
	return a.username
}

func (a *API) ListChannels(teamName string) ([]string, error) {
	apiInput := fmt.Sprintf(`{"method": "listconvsonname", "params": {"options": {"topic_type": "CHAT", "members_type": "team", "name": "%s"}}}`, teamName)
	output, err := a.doFetch(apiInput)
	if err != nil {
		return nil, err
	}

	var channelsList ChannelsList
	if err := json.Unmarshal(output, &channelsList); err != nil {
		return nil, err
	} else if channelsList.Error != nil {
		return nil, errors.New(channelsList.Error.Message)
	}

	var channels []string
	for _, conv := range channelsList.Result.Convs {
		channels = append(channels, conv.Channel.TopicName)
	}
	return channels, nil
}

func (a *API) JoinChannel(teamName string, channelName string) (chat1.EmptyRes, error) {
	empty := chat1.EmptyRes{}

	apiInput := fmt.Sprintf(`{"method": "join", "params": {"options": {"channel": {"name": "%s", "members_type": "team", "topic_name": "%s"}}}}`, teamName, channelName)
	output, err := a.doFetch(apiInput)
	if err != nil {
		return empty, err
	}

	joinChannel := JoinChannel{}
	err = json.Unmarshal(output, &joinChannel)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase team api: %v", err)
	} else if joinChannel.Error != nil {
		return empty, errors.New(joinChannel.Error.Message)
	}

	return joinChannel.Result, nil
}

func (a *API) LeaveChannel(teamName string, channelName string) (chat1.EmptyRes, error) {
	empty := chat1.EmptyRes{}

	apiInput := fmt.Sprintf(`{"method": "leave", "params": {"options": {"channel": {"name": "%s", "members_type": "team", "topic_name": "%s"}}}}`, teamName, channelName)
	output, err := a.doFetch(apiInput)
	if err != nil {
		return empty, err
	}

	leaveChannel := LeaveChannel{}
	err = json.Unmarshal(output, &leaveChannel)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase team api: %v", err)
	} else if leaveChannel.Error != nil {
		return empty, errors.New(leaveChannel.Error.Message)
	}

	return leaveChannel.Result, nil
}

func (a *API) LogSend(feedback string) error {
	feedback = "go-keybase-chat-bot log send\n" +
		"username: " + a.GetUsername() + "\n" +
		feedback

	args := []string{
		"log", "send",
		"--no-confirm",
		"--feedback", feedback,
	}

	// We're determining whether the service is already running by running status
	// with autofork disabled.
	if err := a.runOpts.Command("--no-auto-fork", "status"); err != nil {
		// Assume that there's no service running, so log send as standalone
		args = append([]string{"--standalone"}, args...)
	}

	return a.runOpts.Command(args...).Run()
}

func (a *API) Shutdown() error {
	if a.runOpts.Oneshot != nil {
		err := a.runOpts.Command("logout", "--force").Run()
		if err != nil {
			return err
		}
	}

	if a.runOpts.StartService {
		err := a.runOpts.Command("ctl", "stop", "--shutdown").Run()
		if err != nil {
			return err
		}
	}

	return nil
}
