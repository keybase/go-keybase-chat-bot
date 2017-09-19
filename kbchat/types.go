package kbchat

type Sender struct {
	Uid        string `json:"uid"`
	Username   string `json:"username"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:device_name"`
}

type Channel struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	TopicType   string `json:"topic_type"`
	TopicName   string `json:"topic_name"`
	MembersType string `json:"members_type"`
}

type Conversation struct {
	Id      string  `json:"id"`
	Unread  bool    `json:"unread"`
	Channel Channel `json:"channel"`
}

type Result struct {
	Convs []Conversation `json:"conversations"`
}

type Inbox struct {
	Result Result `json:"result"`
}

type Text struct {
	Body string `json:"body"`
}

type Content struct {
	Type string `json:"type"`
	Text Text   `json:"text"`
}

type Message struct {
	Content Content `json:"content"`
	Sender  Sender  `json:"sender"`
}

type MessageHolder struct {
	Msg Message `json:"msg"`
}

type ThreadResult struct {
	Messages []MessageHolder `json:"messages"`
}

type Thread struct {
	Result ThreadResult `json:"result"`
}
