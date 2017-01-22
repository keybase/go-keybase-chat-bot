package kbchat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type API struct {
	input  io.Writer
	output *bufio.Scanner
}

func NewAPI(keybaseLocation string) (*API, error) {

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
		input:  input,
		output: boutput,
	}, nil
}

func (a *API) GetConversations(unreadOnly bool) ([]Conversation, error) {
	list := fmt.Sprintf("{\"method\":\"list\", \"params\": { \"options\": { \"unread_only\": %v}}}", unreadOnly)
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

func (a *API) GetTextMessages(convID string, unreadOnly bool) ([]Message, error) {
	read := fmt.Sprintf("{\"method\": \"read\", \"params\": {\"options\": {\"conversation_id\": \"%s\", \"unread_only\": %v}}}", convID, unreadOnly)
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

func (a *API) SendMessage(convID string, body string) error {
	send := fmt.Sprintf("{\"method\": \"send\", \"params\": {\"options\": {\"conversation_id\": \"%s\", \"message\": {\"body\": \"%s\"}}}}", convID, body)
	if _, err := io.WriteString(a.input, send); err != nil {
		return err
	}
	a.output.Scan()
	return nil
}

type SubscriptionMessage struct {
	Message      Message
	Conversation Conversation
}

type NewMessageSubscription struct {
	newMsgsCh  <-chan SubscriptionMessage
	errorCh    <-chan error
	shutdownCh chan struct{}
}

func (m NewMessageSubscription) Read() (SubscriptionMessage, error) {
	select {
	case msg := <-m.newMsgsCh:
		return msg, nil
	case err := <-m.errorCh:
		return SubscriptionMessage{}, err
	}
}

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
