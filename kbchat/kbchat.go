package kbchat

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// API is the main object used for communicating with the Keybase JSON API
type API struct {
	input    io.Writer
	output   *bufio.Scanner
	username string
	runOpts  RunOptions
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

	// If a paper key is specified, then login with oneshot mode (logout first)
	if runOpts.Oneshot != nil {
		if err := runOpts.Command("logout", "-f").Run(); err != nil {
			return nil, err
		}
		if err := runOpts.Command("oneshot", "--username", runOpts.Oneshot.Username, "--paperkey",
			runOpts.Oneshot.PaperKey).Run(); err != nil {
			return nil, err
		}
	}

	// Get username first
	username, err := getUsername(runOpts)
	if err != nil {
		return nil, err
	}

	p := runOpts.Command("chat", "api")

	input, err := p.StdinPipe()
	if err != nil {
		return nil, err
	}
	output, err := p.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := p.Start(); err != nil {
		return nil, err
	}

	boutput := bufio.NewScanner(output)
	return &API{
		input:    input,
		output:   boutput,
		username: username,
		runOpts:  runOpts,
	}, nil
}

// GetConversations reads all conversations from the current user's inbox.
func (a *API) GetConversations(unreadOnly bool) ([]Conversation, error) {
	list := fmt.Sprintf(`{"method":"list", "params": { "options": { "unread_only": %v}}}`, unreadOnly)
	if _, err := io.WriteString(a.input, list); err != nil {
		return nil, err
	}
	a.output.Scan()

	var inbox Inbox
	inboxRaw := a.output.Text()
	if err := json.Unmarshal([]byte(inboxRaw[:]), &inbox); err != nil {
		return nil, err
	}
	return inbox.Result.Convs, nil
}

// GetTextMessages fetches all text messages from a given channel. Optionally can filter
// ont unread status.
func (a *API) GetTextMessages(channel Channel, unreadOnly bool) ([]Message, error) {
	channelBytes, err := json.Marshal(channel)
	if err != nil {
		return nil, err
	}

	read := fmt.Sprintf(`{"method": "read", "params": {"options": {"channel": %s}}}`, string(channelBytes))
	if _, err := io.WriteString(a.input, read); err != nil {
		return nil, err
	}
	a.output.Scan()

	var thread Thread
	if err := json.Unmarshal([]byte(a.output.Text()), &thread); err != nil {
		return nil, fmt.Errorf("unable to decode thread: %s", err.Error())
	}

	var res []Message
	for _, msg := range thread.Result.Messages {
		if msg.Msg.Content.Type == "text" {
			res = append(res, msg.Msg)
		}
	}

	return res, nil
}

type sendMessageBody struct {
	Body string
}

type sendMessageOptions struct {
	Channel        Channel         `json:"channel,omitempty"`
	ConversationID string          `json:"conversation_id,omitempty"`
	Message        sendMessageBody `json:",omitempty"`
	Filename       string          `json:"filename,omitempty"`
	Title          string          `json:"title,omitempty"`
}

type sendMessageParams struct {
	Options sendMessageOptions
}

type sendMessageArg struct {
	Method string
	Params sendMessageParams
}

func (a *API) doSend(arg sendMessageArg) error {
	bArg, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(a.input, string(bArg)); err != nil {
		return err
	}
	a.output.Scan()
	return nil
}

// SendMessage sends a new text message on the given channel
func (a *API) SendMessage(channel Channel, body string) error {
	arg := sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				Channel: channel,
				Message: sendMessageBody{
					Body: body,
				},
			},
		},
	}
	return a.doSend(arg)
}

func (a *API) SendMessageByConvID(convID string, body string) error {
	arg := sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				ConversationID: convID,
				Message: sendMessageBody{
					Body: body,
				},
			},
		},
	}
	return a.doSend(arg)
}

// SendMessageByTlfName sends a message on the given TLF name
func (a *API) SendMessageByTlfName(tlfName string, body string) error {
	arg := sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				Channel: Channel{
					Name: tlfName,
				},
				Message: sendMessageBody{
					Body: body,
				},
			},
		},
	}
	return a.doSend(arg)
}

func (a *API) SendMessageByTeamName(teamName string, body string, inChannel *string) error {
	channel := "general"
	if inChannel != nil {
		channel = *inChannel
	}
	arg := sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				Channel: Channel{
					MembersType: "team",
					Name:        teamName,
					TopicName:   channel,
				},
				Message: sendMessageBody{
					Body: body,
				},
			},
		},
	}
	return a.doSend(arg)
}

func (a *API) SendAttachmentByTeam(teamName string, filename string, title string, inChannel *string) error {
	channel := "general"
	if inChannel != nil {
		channel = *inChannel
	}
	arg := sendMessageArg{
		Method: "attach",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				Channel: Channel{
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

func (a *API) Username() string {
	return a.username
}

// SubscriptionMessage contains a message and conversation object
type SubscriptionMessage struct {
	Message      Message
	Conversation Conversation
}

// NewMessageSubscription has methods to control the background message fetcher loop
type NewMessageSubscription struct {
	newMsgsCh  <-chan SubscriptionMessage
	errorCh    <-chan error
	shutdownCh chan struct{}
}

// Read blocks until a new message arrives
func (m NewMessageSubscription) Read() (SubscriptionMessage, error) {
	select {
	case msg := <-m.newMsgsCh:
		return msg, nil
	case err := <-m.errorCh:
		return SubscriptionMessage{}, err
	}
}

// Shutdown terminates the background process
func (m NewMessageSubscription) Shutdown() {
	m.shutdownCh <- struct{}{}
}

// ListenForNewTextMessages fires off a background loop to fetch incoming messages.
func (a *API) ListenForNewTextMessages() (NewMessageSubscription, error) {
	p := a.runOpts.Command("chat", "api-listen")
	output, err := p.StdoutPipe()
	if err != nil {
		return NewMessageSubscription{}, fmt.Errorf("Failed to listen: %s", err)
	}

	newMsgCh := make(chan SubscriptionMessage, 100)
	errorCh := make(chan error, 100)
	shutdownCh := make(chan struct{})

	sub := NewMessageSubscription{
		newMsgsCh:  newMsgCh,
		shutdownCh: shutdownCh,
		errorCh:    errorCh,
	}

	boutput := bufio.NewScanner(output)
	go func() {
		for {
			select {
			case <-shutdownCh:
				return
			default:
				boutput.Scan()
				t := boutput.Text()
				var holder MessageHolder
				var subscriptionMessage SubscriptionMessage
				if err := json.Unmarshal([]byte(t), &holder); err != nil {
					errorCh <- err
					continue
				}
				subscriptionMessage = SubscriptionMessage{
					Message: holder.Msg,
					Conversation: Conversation{
						Channel: holder.Msg.Channel,
					},
				}
				newMsgCh <- subscriptionMessage
			}
		}
	}()

	if err := p.Start(); err != nil {
		return NewMessageSubscription{}, err
	}

	return sub, nil
}

func (a *API) GetUsername() string {
	return a.username
}
