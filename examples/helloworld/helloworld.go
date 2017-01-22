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

	if kbc, err = kbchat.Start(kbLoc); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	tlfName := fmt.Sprintf("%s,%s", kbc.Username(), "kb_monbot")
	fmt.Printf("saying hello on conversation: %s\n", tlfName)
	if err = kbc.SendMessageByTlfName(tlfName, "hello!"); err != nil {
		fail("Error sending message; %s", err.Error())
	}

}
