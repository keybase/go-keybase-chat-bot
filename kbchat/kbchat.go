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
)

// API is the main object used for communicating with the Keybase JSON API
type API struct {
	sync.Mutex
	apiInput  io.Writer
	apiOutput *bufio.Scanner
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
	} else {
		if a.runOpts.Oneshot == nil {
			return "", err
		}
		username = ""
	}
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
		return username, nil
	}
	return "", errors.New("unable to auth")
}

func (a *API) startPipes() (err error) {
	a.Lock()
	defer a.Unlock()
	if a.apiCmd != nil {
		a.apiCmd.Process.Kill()
	}
	a.apiCmd = nil
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
	a.apiOutput = bufio.NewScanner(output)
	return nil
}

var errAPIDisconnected = errors.New("chat API disconnected")

func (a *API) getAPIPipes() (io.Writer, *bufio.Scanner, error) {
	a.Lock()
	defer a.Unlock()
	if a.apiCmd == nil {
		return nil, nil, errAPIDisconnected
	}
	return a.apiInput, a.apiOutput, nil
}

// GetConversations reads all conversations from the current user's inbox.
func (a *API) GetConversations(unreadOnly bool) ([]Conversation, error) {
	input, output, err := a.getAPIPipes()
	if err != nil {
		return nil, err
	}
	list := fmt.Sprintf(`{"method":"list", "params": { "options": { "unread_only": %v}}}`, unreadOnly)
	if _, err := io.WriteString(input, list); err != nil {
		return nil, err
	}
	output.Scan()

	var inbox Inbox
	inboxRaw := output.Text()
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

	input, output, err := a.getAPIPipes()
	if err != nil {
		return nil, err
	}
	read := fmt.Sprintf(`{"method": "read", "params": {"options": {"channel": %s}}}`, string(channelBytes))
	if _, err := io.WriteString(input, read); err != nil {
		return nil, err
	}
	output.Scan()

	var thread Thread
	if err := json.Unmarshal([]byte(output.Text()), &thread); err != nil {
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
	input, output, err := a.getAPIPipes()
	if err != nil {
		return err
	}
	if _, err := io.WriteString(input, string(bArg)); err != nil {
		return err
	}
	output.Scan()
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

type SubscriptionWalletEvent struct {
	Payment Payment
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
				return
			}
			switch typeHolder.Type {
			case "chat":
				var holder MessageHolder
				if err := json.Unmarshal([]byte(t), &holder); err != nil {
					errorCh <- err
					return
				}
				subscriptionMessage := SubscriptionMessage{
					Message: holder.Msg,
					Conversation: Conversation{
						Channel: holder.Msg.Channel,
					},
				}
				newMsgCh <- subscriptionMessage
			case "wallet":
				var holder PaymentHolder
				if err := json.Unmarshal([]byte(t), &holder); err != nil {
					errorCh <- err
					return
				}
				subscriptionPayment := SubscriptionWalletEvent{
					Payment: holder.Payment,
				}
				newWalletCh <- subscriptionPayment
			default:
				continue
			}
		}
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
			p.Wait()
			time.Sleep(pause)
		}
	}()
	return sub, nil
}

func (a *API) GetUsername() string {
	return a.username
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
