# go-keybase-chat-bot

Script Keybase Chat in Go!

This module is a side-project/work in progress and may change or have crashers, but feel free to play with it. As long as you're logged in as a Keybase user, you can use this module to script basic chat commands.

# Installation

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

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc)}); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	tlfName := fmt.Sprintf("%s,%s", kbc.GetUsername(), "kb_monbot")
	fmt.Printf("saying hello on conversation: %s\n", tlfName)
	if err = kbc.SendMessageByTlfName(tlfName, "hello!"); err != nil {
		fail("Error sending message; %s", err.Error())
	}
}

```

### Commands

#### `Start(keybaseLocation string) *API`

This must be run first in order to start the Keybase JSON API stdin/stdout interactive mode.

#### `API.SendMessageByTlfName(tlfName string, body string) error`

send a new message by specifying a TLF name

#### `API.SendMessage(convID string, body string) error`

send a new message by specifying a conversation ID

#### `API.GetConversations(unreadOnly bool) ([]Conversation, error)`

get all conversations, optionally filtering for unread status

#### `API.GetTextMessages(convID string, unreadOnly bool) ([]Message, error)`

get all text messages, optionally filtering for unread status

Reads the messages in a channel. You can read with or without marking as read.

#### `API.ListenForNextTextMessages() NewSubscription`

Returns an object that allows for a bot to enter into a loop calling `NewSubscription.Read`
to receive any new message across all conversations (except the bots own messages). See the following example:

```go
	sub := kbc.ListenForNewTextMessages()
	for {
		msg, err := sub.Read()
		if err != nil {
			fail("failed to read message: %s", err.Error())
		}

		if err = kbc.SendMessage(msg.Conversation.Id, msg.Message.Content.Text.Body); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}
	}
```

#### `API.Listen(kbc.ListenOptions{Wallet: true}) NewSubscription`

Returns the same object as above, but this one will have another channel on it that also gets wallet events. You can get those just like chat messages: `NewSubscription.ReadWallet`. So if you care about both of these types of events, you might run two loops like this:

```go
	sub := kbc.Listen(kbc.ListenOptions{Wallet: true})
	go func() {
		for {
			payment, err := sub.ReadWallet()
			if err != nil {
				fail("failed to read payment event: %s", err.Error())
			}
			tlfName := fmt.Sprintf("%s,%s", payment.FromUsername, "kb_monbot")
			msg := fmt.Sprintf("thanks for the %s!", payment.AmountDescription)
			if err = kbc.SendMessageByTlfName(tlfName, msg); err != nil {
				fail("error thanking for payment: %s", err.Error())
			}
		}
	}()

	for {
		msg, err := sub.Read()
		if err != nil {
			fail("failed to read message: %s", err.Error())
		}

		if err = kbc.SendMessage(msg.Conversation.Id, msg.Message.Content.Text.Body); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}
	}
```



## TODO:
  - attachment handling (posting/getting)
  - edit/delete
  - many other things!

### Contributions

- welcomed!
