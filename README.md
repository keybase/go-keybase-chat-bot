# go-keybase-chat-bot

[![Travis CI](https://travis-ci.org/keybase/go-keybase-chat-bot.svg?branch=master)](https://travis-ci.org/keybase/go-keybase-chat-bot)

Script Keybase Chat in Go!

This module is a side-project/work in progress and may change or have crashers, but feel free to play with it. As long as you're logged in as a Keybase user, you can use this module to script basic chat commands.

## Installation

Make sure to [install Keybase](https://keybase.io/download).

```bash
go get -u github.com/keybase/go-keybase-chat-bot/...
```

### Hello world

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func main() {
	var kbLoc string
	var kbc *kbchat.API
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	tlfName := fmt.Sprintf("%s,%s", kbc.Username(), "kb_monbot")
	fmt.Printf("saying hello on conversation: %s\n", tlfName)
	if _, err = kbc.SendMessageByTlfName(tlfName, "hello!"); err != nil {
		fail("Error sending message; %s", err.Error())
	}
}

```

### Commands

#### `Start(runOpts RunOptions) (*API, error)`

This must be run first in order to start the Keybase JSON API stdin/stdout interactive mode.

#### `API.SendMessageByTlfName(tlfName string, body string) (SendResponse, error)`

send a new message by specifying a TLF name

#### `API.SendMessage(channel Channel, body string) (SendResponse, error)`

send a new message by specifying a channel

#### `API.SendMessageByConvID(convID string, body string) (SendResponse, error)`

send a new message by specifying a conversation ID

#### `API.GetConversations(unreadOnly bool) ([]Conversation, error)`

get all conversations, optionally filtering for unread status

#### `API.GetTextMessages(channel Channel, unreadOnly bool) ([]Message, error)`

get all text messages, optionally filtering for unread status

Reads the messages in a channel. You can read with or without marking as read.

#### `API.ListenForNextTextMessages() NewSubscription`

Returns an object that allows for a bot to enter into a loop calling `NewSubscription.Read`
to receive any new message across all conversations (except the bots own messages). See the following example:

```go
	sub, err := kbc.ListenForNewTextMessages()
	if err != nil {
		fail("Error listening: %s", err.Error())
	}

	for {
		msg, err := sub.Read()
		if err != nil {
			fail("failed to read message: %s", err.Error())
		}

		if msg.Message.Content.Type != "text" {
			continue
		}

		if msg.Message.Sender.Username == kbc.GetUsername() {
			continue
		}

		if _, err = kbc.SendMessage(msg.Message.Channel, msg.Message.Content.Text.Body); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}
	}
```

#### `API.Listen(kbchat.ListenOptions{Wallet: true}) NewSubscription`

Returns the same object as above, but this one will have another channel on it that also gets wallet events. You can get those just like chat messages: `NewSubscription.ReadWallet`. So if you care about both of these types of events, you might run two loops like this:

```go
	sub, err := kbc.Listen(kbchat.ListenOptions{Wallet: true})
	if err != nil {
		fail("Error listening: %s", err.Error())
	}

	go func() {
		for {
			payment, err := sub.ReadWallet()
			if err != nil {
				fail("failed to read payment event: %s", err.Error())
			}
			tlfName := fmt.Sprintf("%s,%s", payment.Payment.FromUsername, "kb_monbot")
			msg := fmt.Sprintf("thanks for the %s!", payment.Payment.AmountDescription)
			if _, err = kbc.SendMessageByTlfName(tlfName, msg); err != nil {
				fail("error thanking for payment: %s", err.Error())
			}
		}
	}()

	for {
		msg, err := sub.Read()
		if err != nil {
			fail("failed to read message: %s", err.Error())
		}

		if msg.Message.Content.Type != "text" {
			continue
		}

		if msg.Message.Sender.Username == kbc.GetUsername() {
			continue
		}

		if _, err = kbc.SendMessage(msg.Message.Channel, msg.Message.Content.Text.Body); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}
	}
```

## TODO:

- attachment handling (posting/getting)
- edit/delete
- many other things!

## Contributions

- welcomed!

### Precommit hooks

We check all git commits with pre-commit hooks generated via
[pre-commit.com](http://pre-commit.com) pre-commit hooks.
To enable use of these pre-commit hooks:

- [Install](http://pre-commit.com/#install) the `pre-commit` utility. For some common cases:
  - `pip install pre-commit`
  - `brew install pre-commit`
- Remove any existing pre-commit hooks via `rm .git/hooks/pre-commit`
- Configure via `pre-commit install`

### Testing

You'll need to have a few demo bot accounts and teams to run the suite of tests. Make a copy of `kbchat/test_config.example.yaml`, rename it to `kbchat/test_config.yaml`, and replace the example data with your own. Tests can then be run inside of `kbchat/` with `go test`.
