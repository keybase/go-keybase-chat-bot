package kbchat

import "github.com/keybase/go-keybase-chat-bot/kbchat/types/chat1"

type PaymentHolder struct {
	Payment Payment `json:"notification"`
}

type Payment struct {
	TxID              string `json:"txID"`
	StatusDescription string `json:"statusDescription"`
	FromAccountID     string `json:"fromAccountID"`
	FromUsername      string `json:"fromUsername"`
	ToAccountID       string `json:"toAccountID"`
	ToUsername        string `json:"toUsername"`
	AmountDescription string `json:"amountDescription"`
	WorthAtSendTime   string `json:"worthAtSendTime"`
	ExternalTxURL     string `json:"externalTxURL"`
}

type Result struct {
	Convs []chat1.ConvSummary `json:"conversations"`
}

type Inbox struct {
	Result Result `json:"result"`
}

type ChannelsList struct {
	Result Result `json:"result"`
}

type SendResponse struct {
	Result chat1.SendRes `json:"result"`
}

type TypeHolder struct {
	Type string `json:"type"`
}

type Thread struct {
	Result chat1.Thread `json:"result"`
}

type Advertisement struct {
	Alias          string `json:"alias,omitempty"`
	Advertisements []chat1.AdvertiseCommandsParam
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type JoinChannel struct {
	Error  Error          `json:"error"`
	Result chat1.EmptyRes `json:"result"`
}

type LeaveChannel struct {
	Error  Error          `json:"error"`
	Result chat1.EmptyRes `json:"result"`
}
