package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
)

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func main() {
	const MsgPrefix = "!storage"
	const (
		Help   = "help"
		List   = "list"
		Get    = "get"
		Put    = "put"
		Delete = "delete"
	)

	var kbLoc string
	var kbc *kbchat.API
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	fmt.Println("Simple kvstore storage bot now starting...")

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	sub, err := kbc.ListenForNewTextMessages()
	if err != nil {
		fail("Error listening: %s", err.Error())
	}

	fmt.Println("Now listening for new chat messages...")
	for {
		msg, err := sub.Read()

		if err != nil {
			fail("failed to read message: %s", err.Error())
		}

		if msg.Message.Content.TypeName != "text" {
			continue
		}

		channel := msg.Message.Channel
		team := channel.Name
		body := strings.Split(strings.TrimSpace(msg.Message.Content.Text.Body), " ")

		if len(body) < 2 || body[0] != MsgPrefix {
			continue
		}

		switch action := body[1]; action {
		case Help:
			handleHelp(kbc, channel, team, body, action)
		case List:
			handleList(kbc, channel, team, body, action)
		case Get:
			handleGet(kbc, channel, team, body, action)
		case Put:
			handlePut(kbc, channel, team, body, action)
		case Delete:
			handleDelete(kbc, channel, team, body, action)
		}

	}

}

func handleHelp(api *kbchat.API, channel chat1.ChatChannel, team string, body []string, action string) {
	// !storage help
	sendMsg :=
		"Available commands:" +
			"\n`!storage put <namespace> <key> <value> (<revision>)`" +
			"\n`!storage get <namespace> <key>`" +
			"\n`!storage delete <namespace> <key> (<revision>)`" +
			"\n`!storage list`  // list namespaces" +
			"\n`!storage list <namespace>`  // list entries in namespace"

	if _, err := api.SendMessage(channel, sendMsg); err != nil {
		fail("error sending message: %s", err.Error())
	}
}

func handleList(api *kbchat.API, channel chat1.ChatChannel, team string, body []string, action string) {
	sendMsg := "Error handling list command"

	if len(body) == 2 {
		// !storage list
		if res, err := api.ListNamespaces(team); err == nil {
			sendMsg = fmt.Sprintf("%+v", res)
		} else {
			sendMsg = fmt.Sprintf("%+v", err)
		}
	} else if len(body) == 3 {
		// !storage list ,namespace.
		namespace := body[2]
		if res, err := api.ListEntryKeys(team, namespace); err == nil {
			sendMsg = fmt.Sprintf("%+v", res)
		} else {
			sendMsg = fmt.Sprintf("%+v", err)
		}
	}

	if _, err := api.SendMessage(channel, sendMsg); err != nil {
		fail("error sending message: %s", err.Error())
	}
	return
}

func handleGet(api *kbchat.API, channel chat1.ChatChannel, team string, body []string, action string) {
	sendMsg := "Error handling get command"

	if len(body) == 4 {
		// !storage get <namespace> <key>
		namespace, key := body[2], body[3]
		if res, err := api.GetEntry(team, namespace, key); err == nil {
			sendMsg = fmt.Sprintf("%+v", res)
		} else {
			sendMsg = fmt.Sprintf("%+v", err)
		}
	}

	if _, err := api.SendMessage(channel, sendMsg); err != nil {
		fail("error sending message: %s", err.Error())
	}
	return
}

func handlePut(api *kbchat.API, channel chat1.ChatChannel, team string, body []string, action string) {
	sendMsg := "Error handling put command"

	if len(body) == 5 || len(body) == 6 {
		// !storage put <namespace> <key> <value> (<revision>)
		namespace, key, value := body[2], body[3], body[4]
		revision := 0
		if len(body) == 6 {
			thisRevision, err := strconv.Atoi(body[5])
			if err != nil {
				if _, err := api.SendMessage(channel, sendMsg); err != nil {
					fail("error sending message: %s", err.Error())
				}
				return
			}
			revision = thisRevision
		}
		if res, err := api.PutEntry(team, namespace, key, value, revision); err == nil {
			sendMsg = fmt.Sprintf("%+v", res)
		} else {
			sendMsg = fmt.Sprintf("%+v", err)
		}
	}

	if _, err := api.SendMessage(channel, sendMsg); err != nil {
		fail("error sending message: %s", err.Error())
	}
}

func handleDelete(api *kbchat.API, channel chat1.ChatChannel, team string, body []string, action string) {
	sendMsg := "Error handling delete command"
	if len(body) == 4 || len(body) == 5 {
		// !storage delete <namespace> <key> (<revision>)
		namespace, key := body[2], body[3]
		revision := 0 ///////
		if len(body) == 5 {
			thisRevision, err := strconv.Atoi(body[4])
			if err != nil {
				if _, err := api.SendMessage(channel, sendMsg); err != nil {
					fail("error sending message: %s", err.Error())
				}
				return
			}
			revision = thisRevision
		}
		if res, err := api.DeleteEntry(team, namespace, key, revision); err == nil {
			sendMsg = fmt.Sprintf("%+v", res)
		} else {
			sendMsg = fmt.Sprintf("%+v", err)
		}

	}
	if _, err := api.SendMessage(channel, sendMsg); err != nil {
		fail("error sending message: %s", err.Error())
	}
}
