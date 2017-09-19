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
}

func getUsername(keybaseLocation string) (username string, err error) {
	p := exec.Command(keybaseLocation, "status")
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

// Start fires up the Keybase JSON API in stdin/stdout mode
func Start(keybaseLocation string) (*API, error) {

	// Get username first
	username, err := getUsername(keybaseLocation)
	if err != nil {
		return nil, err
	}

	p := exec.Command(keybaseLocation, "chat", "api")
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
	}, nil
}

// GetConversations reads all conversations from the current user's inbox. Optionally
// can filter for unread only.
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

// GetTextMessages fetches all text messages from a given conversation ID. Optionally can filter
// ont unread status.
func (a *API) GetTextMessages(convID string, unreadOnly bool) ([]Message, error) {
	read := fmt.Sprintf(`{"method": "read", "params": {"options": {"conversation_id": "%s", "unread_only": %v}}}`, convID, unreadOnly)
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
	ConversationID string  `json:"conversation_id,omitempty"`
	Channel        Channel `json:"channel,omitempty"`
	Message        sendMessageBody
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

// SendMessage sends a new text message on the given conversation ID
func (a *API) SendMessage(convID string, body string) error {
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

func (a *API) getUnreadMessagesFromConvs(convs []Conversation) ([]SubscriptionMessage, error) {
	var res []SubscriptionMessage
	for _, conv := range convs {
		msgs, err := a.GetTextMessages(conv.Id, true)
		if err != nil {
			return nil, err
		}
		for _, msg := range msgs {
			res = append(res, SubscriptionMessage{
				Message:      msg,
				Conversation: conv,
			})
		}
	}
	return res, nil
}

// ListenForNewTextMessages fires off a background loop to fetch new unread messages.
func (a *API) ListenForNewTextMessages() NewMessageSubscription {
	newMsgCh := make(chan SubscriptionMessage, 100)
	errorCh := make(chan error, 100)
	shutdownCh := make(chan struct{})
	sub := NewMessageSubscription{
		newMsgsCh:  newMsgCh,
		shutdownCh: shutdownCh,
		errorCh:    errorCh,
	}
	go func() {
		for {
			select {
			case <-shutdownCh:
				return
			case <-time.After(2 * time.Second):
				// Get all unread convos
				convs, err := a.GetConversations(true)
				if err != nil {
					errorCh <- err
					continue
				}
				// Get unread msgs from convs
				msgs, err := a.getUnreadMessagesFromConvs(convs)
				if err != nil {
					errorCh <- err
					continue
				}
				// Send all the new messages out
				for _, msg := range msgs {
					newMsgCh <- msg
				}
			}
		}
	}()

	return sub
}
