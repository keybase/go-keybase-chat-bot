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

	sub, err := kbc.ListenForNewTextMessages()
	if err != nil {
		fail("Error listening: %s", err.Error())
	}

	for {
		msg, err := sub.Read()
		if err != nil {
			fail("failed to read message: %s", err.Error())
		}

		if !(msg.Message.Sender.Username == "modalduality" && msg.Message.Channel.Name == "scianmuses") {
			continue
		}

		if msg.Message.Content.Type != "text" {
			continue
		}

		if err = kbc.SendMessage(msg.Message.Channel, msg.Message.Content.Text.Body); err != nil {
			fail("error echo'ing message: %s", err.Error())
		}
	}

}
