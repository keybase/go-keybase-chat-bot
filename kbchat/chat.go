package kbchat

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
)

type Thread struct {
	Result chat1.Thread `json:"result"`
	Error  *Error       `json:"error,omitempty"`
}

type Inbox struct {
	Result Result `json:"result"`
	Error  *Error `json:"error,omitempty"`
}

type sendMessageBody struct {
	Body string
}

type sendMessageOptions struct {
	Channel          chat1.ChatChannel `json:"channel,omitempty"`
	ConversationID   string            `json:"conversation_id,omitempty"`
	Message          sendMessageBody   `json:",omitempty"`
	Filename         string            `json:"filename,omitempty"`
	Title            string            `json:"title,omitempty"`
	MsgID            chat1.MessageID   `json:"message_id,omitempty"`
	ConfirmLumenSend bool              `json:"confirm_lumen_send"`
}

type sendMessageParams struct {
	Options sendMessageOptions
}

type sendMessageArg struct {
	Method string
	Params sendMessageParams
}

func newSendArg(options sendMessageOptions) sendMessageArg {
	return sendMessageArg{
		Method: "send",
		Params: sendMessageParams{
			Options: options,
		},
	}
}

// GetConversations reads all conversations from the current user's inbox.
func (a *API) GetConversations(unreadOnly bool) ([]chat1.ConvSummary, error) {
	apiInput := fmt.Sprintf(`{"method":"list", "params": { "options": { "unread_only": %v}}}`, unreadOnly)
	output, err := a.doFetch(apiInput)
	if err != nil {
		return nil, err
	}

	var inbox Inbox
	if err := json.Unmarshal(output, &inbox); err != nil {
		return nil, err
	} else if inbox.Error != nil {
		return nil, errors.New(inbox.Error.Message)
	}
	return inbox.Result.Convs, nil
}

// GetTextMessages fetches all text messages from a given channel. Optionally can filter
// ont unread status.
func (a *API) GetTextMessages(channel chat1.ChatChannel, unreadOnly bool) ([]chat1.MsgSummary, error) {
	channelBytes, err := json.Marshal(channel)
	if err != nil {
		return nil, err
	}
	apiInput := fmt.Sprintf(`{"method": "read", "params": {"options": {"channel": %s}}}`, string(channelBytes))
	output, err := a.doFetch(apiInput)
	if err != nil {
		return nil, err
	}

	var thread Thread

	if err := json.Unmarshal(output, &thread); err != nil {
		return nil, fmt.Errorf("unable to decode thread: %v", err)
	} else if thread.Error != nil {
		return nil, errors.New(thread.Error.Message)
	}

	var res []chat1.MsgSummary
	for _, msg := range thread.Result.Messages {
		if msg.Msg.Content.TypeName == "text" {
			res = append(res, *msg.Msg)
		}
	}

	return res, nil
}

func (a *API) SendMessage(channel chat1.ChatChannel, body string, args ...interface{}) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		Channel: channel,
		Message: sendMessageBody{
			Body: fmt.Sprintf(body, args...),
		},
	})
	return a.doSend(arg)
}

func (a *API) Broadcast(body string, args ...interface{}) (SendResponse, error) {
	return a.SendMessage(chat1.ChatChannel{
		Name:   a.GetUsername(),
		Public: true,
	}, fmt.Sprintf(body, args...))
}

func (a *API) SendMessageByConvID(convID string, body string, args ...interface{}) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		ConversationID: convID,
		Message: sendMessageBody{
			Body: fmt.Sprintf(body, args...),
		},
	})
	return a.doSend(arg)
}

// SendMessageByTlfName sends a message on the given TLF name
func (a *API) SendMessageByTlfName(tlfName string, body string, args ...interface{}) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		Channel: chat1.ChatChannel{
			Name: tlfName,
		},
		Message: sendMessageBody{
			Body: fmt.Sprintf(body, args...),
		},
	})
	return a.doSend(arg)
}

func (a *API) SendMessageByTeamName(teamName string, inChannel *string, body string, args ...interface{}) (SendResponse, error) {
	channel := "general"
	if inChannel != nil {
		channel = *inChannel
	}
	arg := newSendArg(sendMessageOptions{
		Channel: chat1.ChatChannel{
			MembersType: "team",
			Name:        teamName,
			TopicName:   channel,
		},
		Message: sendMessageBody{
			Body: fmt.Sprintf(body, args...),
		},
	})
	return a.doSend(arg)
}

func (a *API) SendAttachmentByTeam(teamName string, inChannel *string, filename string, title string) (SendResponse, error) {
	channel := "general"
	if inChannel != nil {
		channel = *inChannel
	}
	arg := sendMessageArg{
		Method: "attach",
		Params: sendMessageParams{
			Options: sendMessageOptions{
				Channel: chat1.ChatChannel{
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

////////////////////////////////////////////////////////
// React to chat ///////////////////////////////////////
////////////////////////////////////////////////////////

type reactionOptions struct {
	ConversationID string `json:"conversation_id"`
	Message        sendMessageBody
	MsgID          chat1.MessageID   `json:"message_id"`
	Channel        chat1.ChatChannel `json:"channel"`
}

type reactionParams struct {
	Options reactionOptions
}

type reactionArg struct {
	Method string
	Params reactionParams
}

func newReactionArg(options reactionOptions) reactionArg {
	return reactionArg{
		Method: "reaction",
		Params: reactionParams{Options: options},
	}
}

func (a *API) ReactByChannel(channel chat1.ChatChannel, msgID chat1.MessageID, reaction string) (SendResponse, error) {
	arg := newReactionArg(reactionOptions{
		Message: sendMessageBody{Body: reaction},
		MsgID:   msgID,
		Channel: channel,
	})
	return a.doSend(arg)
}

func (a *API) ReactByConvID(convID string, msgID chat1.MessageID, reaction string) (SendResponse, error) {
	arg := newReactionArg(reactionOptions{
		Message:        sendMessageBody{Body: reaction},
		MsgID:          msgID,
		ConversationID: convID,
	})
	return a.doSend(arg)
}

////////////////////////////////////////////////////////
// Manage channels /////////////////////////////////////
////////////////////////////////////////////////////////

type ChannelsList struct {
	Result Result `json:"result"`
	Error  *Error `json:"error,omitempty"`
}

type JoinChannel struct {
	Error  *Error         `json:"error,omitempty"`
	Result chat1.EmptyRes `json:"result"`
}

type LeaveChannel struct {
	Error  *Error         `json:"error,omitempty"`
	Result chat1.EmptyRes `json:"result"`
}

func (a *API) ListChannels(teamName string) ([]string, error) {
	apiInput := fmt.Sprintf(`{"method": "listconvsonname", "params": {"options": {"topic_type": "CHAT", "members_type": "team", "name": "%s"}}}`, teamName)
	output, err := a.doFetch(apiInput)
	if err != nil {
		return nil, err
	}

	var channelsList ChannelsList
	if err := json.Unmarshal(output, &channelsList); err != nil {
		return nil, err
	} else if channelsList.Error != nil {
		return nil, errors.New(channelsList.Error.Message)
	}

	var channels []string
	for _, conv := range channelsList.Result.Convs {
		channels = append(channels, conv.Channel.TopicName)
	}
	return channels, nil
}

func (a *API) JoinChannel(teamName string, channelName string) (chat1.EmptyRes, error) {
	empty := chat1.EmptyRes{}

	apiInput := fmt.Sprintf(`{"method": "join", "params": {"options": {"channel": {"name": "%s", "members_type": "team", "topic_name": "%s"}}}}`, teamName, channelName)
	output, err := a.doFetch(apiInput)
	if err != nil {
		return empty, err
	}

	res := JoinChannel{}
	err = json.Unmarshal(output, &res)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase team api: %v", err)
	} else if res.Error != nil {
		return empty, errors.New(res.Error.Message)
	}

	return res.Result, nil
}

func (a *API) LeaveChannel(teamName string, channelName string) (chat1.EmptyRes, error) {
	empty := chat1.EmptyRes{}

	apiInput := fmt.Sprintf(`{"method": "leave", "params": {"options": {"channel": {"name": "%s", "members_type": "team", "topic_name": "%s"}}}}`, teamName, channelName)
	output, err := a.doFetch(apiInput)
	if err != nil {
		return empty, err
	}

	res := LeaveChannel{}
	err = json.Unmarshal(output, &res)
	if err != nil {
		return empty, fmt.Errorf("failed to parse output from keybase team api: %v", err)
	} else if res.Error != nil {
		return empty, errors.New(res.Error.Message)
	}

	return res.Result, nil
}

////////////////////////////////////////////////////////
// Send lumens in chat /////////////////////////////////
////////////////////////////////////////////////////////

func (a *API) InChatSend(channel chat1.ChatChannel, body string, args ...interface{}) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		Channel: channel,
		Message: sendMessageBody{
			Body: fmt.Sprintf(body, args...),
		},
		ConfirmLumenSend: true,
	})
	return a.doSend(arg)
}

func (a *API) InChatSendByConvID(convID string, body string, args ...interface{}) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		ConversationID: convID,
		Message: sendMessageBody{
			Body: fmt.Sprintf(body, args...),
		},
		ConfirmLumenSend: true,
	})
	return a.doSend(arg)
}

func (a *API) InChatSendByTlfName(tlfName string, body string, args ...interface{}) (SendResponse, error) {
	arg := newSendArg(sendMessageOptions{
		Channel: chat1.ChatChannel{
			Name: tlfName,
		},
		Message: sendMessageBody{
			Body: fmt.Sprintf(body, args...),
		},
		ConfirmLumenSend: true,
	})
	return a.doSend(arg)
}

////////////////////////////////////////////////////////
// Misc commands ///////////////////////////////////////
////////////////////////////////////////////////////////

type Advertisement struct {
	Alias          string `json:"alias,omitempty"`
	Advertisements []chat1.AdvertiseCommandAPIParam
}

type ListCommandsResponse struct {
	Result struct {
		Commands []chat1.UserBotCommandOutput `json:"commands"`
	} `json:"result"`
	Error *Error `json:"error,omitempty"`
}

type advertiseCmdsParams struct {
	Options Advertisement
}

type advertiseCmdsMsgArg struct {
	Method string
	Params advertiseCmdsParams
}

func newAdvertiseCmdsMsgArg(ad Advertisement) advertiseCmdsMsgArg {
	return advertiseCmdsMsgArg{
		Method: "advertisecommands",
		Params: advertiseCmdsParams{
			Options: ad,
		},
	}
}

func (a *API) AdvertiseCommands(ad Advertisement) (SendResponse, error) {
	return a.doSend(newAdvertiseCmdsMsgArg(ad))
}

func (a *API) ClearCommands() error {
	arg := struct {
		Method string
	}{
		Method: "clearcommands",
	}
	_, err := a.doSend(arg)
	return err
}

type listCmdsOptions struct {
	Channel        chat1.ChatChannel
	ConversationID string
}

type listCmdsParams struct {
	Options listCmdsOptions
}

type listCmdsArg struct {
	Method string
	Params listCmdsParams
}

func newListCmdsArg(options listCmdsOptions) listCmdsArg {
	return listCmdsArg{
		Method: "listcommands",
		Params: listCmdsParams{
			Options: options,
		},
	}
}

func (a *API) ListCommands(channel chat1.ChatChannel) ([]chat1.UserBotCommandOutput, error) {
	arg := newListCmdsArg(listCmdsOptions{
		Channel: channel,
	})
	return a.listCommands(arg)
}

func (a *API) ListCommandsByConvID(conversationID string) ([]chat1.UserBotCommandOutput, error) {
	arg := newListCmdsArg(listCmdsOptions{
		ConversationID: conversationID,
	})
	return a.listCommands(arg)
}

func (a *API) listCommands(arg listCmdsArg) ([]chat1.UserBotCommandOutput, error) {
	bArg, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}
	output, err := a.doFetch(string(bArg))
	if err != nil {
		return nil, err
	}
	var res ListCommandsResponse
	if err := json.Unmarshal(output, &res); err != nil {
		return nil, err
	} else if res.Error != nil {
		return nil, errors.New(res.Error.Message)
	}
	return res.Result.Commands, nil
}
