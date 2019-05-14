package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

var kbc *kbchat.API

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "fatal error: "+msg+"\n", args...)
	os.Exit(3)
}

func failAndLog(msg string, args ...interface{}) {
	if err := kbc.LogSend(fmt.Sprintf(msg, args...)); err != nil {
		fmt.Fprintf(os.Stderr, "failed to log send: %s\n", err)
	}

	fail(msg, args...)
}

func main() {
	var (
		kbLoc string
		err   error
	)

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	sub, err := kbc.ListenForNewTextMessages()
	if err != nil {
		failAndLog("Error listening: %s", err.Error())
	}

	for {
		msg, err := sub.Read()
		if err != nil {
			failAndLog("failed to read message: %s", err.Error())
		}

		if msg.Message.Content.Type != "text" {
			continue
		}

		if msg.Message.Sender.Username == kbc.GetUsername() {
			continue
		}

		if err = kbc.SendMessage(msg.Message.Channel, msg.Message.Content.Text.Body); err != nil {
			failAndLog("error echo'ing message: %s", err.Error())
		}
	}

}
