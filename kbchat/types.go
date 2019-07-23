package kbchat

import "github.com/keybase/go-keybase-chat-bot/kbchat/types/stellar1"

type Sender struct {
	Uid        string `json:"uid"`
	Username   string `json:"username"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

type Channel struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	TopicType   string `json:"topic_type"`
	TopicName   string `json:"topic_name"`
	MembersType string `json:"members_type"`
}

type Conversation struct {
	ID      string  `json:"id"`
	Unread  bool    `json:"unread"`
	Channel Channel `json:"channel"`
}

type PaymentHolder struct {
	Payment stellar1.PaymentDetailsLocal `json:"notification"`
}

type Result struct {
	Convs []Conversation `json:"conversations"`
}

type Inbox struct {
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

type Message struct {
	Content        Content `json:"content"`
	Sender         Sender  `json:"sender"`
	Channel        Channel `json:"channel"`
	ConversationID string  `json:"conversation_id"`
	MsgID          int     `json:"id"`
}

type SendResult struct {
	MsgID int `json:"id"`
}

type SendResponse struct {
	Result SendResult `json:"result"`
}

type TypeHolder struct {
	Type string `json:"type"`
}

type MessageHolder struct {
	Msg    Message `json:"msg"`
	Source string  `json:"source"`
}

type ThreadResult struct {
	Messages []MessageHolder `json:"messages"`
}

type Thread struct {
	Result ThreadResult `json:"result"`
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
