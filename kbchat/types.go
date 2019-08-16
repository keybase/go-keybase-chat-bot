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

type MsgPaymentDetails struct {
	ResultType int    `json:"resultTyp"` // 0 good. 1 error
	PaymentID  string `json:"sent"`
}

type MsgPayment struct {
	Username    string            `json:"username"`
	PaymentText string            `json:"paymentText"`
	Details     MsgPaymentDetails `json:"result"`
}

type Text struct {
	Body     string       `json:"body"`
	Payments []MsgPayment `json:"payments"`
	ReplyTo  int          `json:"replyTo"`
}

type Content struct {
	Type string `json:"type"`
	Text Text   `json:"text"`
}

type SendResponse struct {
	Result chat1.SendRes `json:"result"`
}

type TypeHolder struct {
	Type string `json:"type"`
}

type MessageHolder struct {
	Msg    chat1.MsgSummary `json:"msg"`
	Source string           `json:"source"`
}

type ThreadResult struct {
	Messages []MessageHolder `json:"messages"`
}
type Thread struct {
	Result chat1.Thread `json:"result"`
}

type CommandExtendedDescription struct {
	Title       string `json:"title"`
	DesktopBody string `json:"desktop_body"`
	MobileBody  string `json:"mobile_body"`
}

type Command struct {
	Name                string                      `json:"name"`
	Description         string                      `json:"description"`
	Usage               string                      `json:"usage"`
	ExtendedDescription *CommandExtendedDescription `json:"extended_description,omitempty"`
}

type CommandsAdvertisement struct {
	Typ      string `json:"type"`
	Commands []Command
	TeamName string `json:"team_name,omitempty"`
}

type Advertisement struct {
	Alias          string `json:"alias,omitempty"`
	Advertisements []CommandsAdvertisement
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
