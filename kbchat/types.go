package kbchat

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
	Unread  bool    `json:"unread"`
	Channel Channel `json:"channel"`
}

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
