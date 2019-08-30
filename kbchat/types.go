package kbchat

import (
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"
	"github.com/keybase/go-keybase-chat-bot/kbchat/types/stellar1"
)

type PaymentHolder struct {
	Payment stellar1.PaymentDetailsLocal `json:"notification"`
}

type Result struct {
	Convs []chat1.ConvSummary `json:"conversations"`
}

type Inbox struct {
	Result Result `json:"result"`
	Error  *Error `json:"error,omitempty"`
}

type ChannelsList struct {
	Result Result `json:"result"`
	Error  *Error `json:"error,omitempty"`
}

type SendResponse struct {
	Result chat1.SendRes `json:"result"`
	Error  *Error        `json:"error,omitempty"`
}

type TypeHolder struct {
	Type string `json:"type"`
}

type Thread struct {
	Result chat1.Thread `json:"result"`
	Error  *Error       `json:"error,omitempty"`
}

type Advertisement struct {
	Alias          string `json:"alias,omitempty"`
	Advertisements []chat1.AdvertiseCommandAPIParam
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JoinChannel struct {
	Error  *Error         `json:"error,omitempty"`
	Result chat1.EmptyRes `json:"result"`
}

type LeaveChannel struct {
	Error  *Error         `json:"error,omitempty"`
	Result chat1.EmptyRes `json:"result"`
}
