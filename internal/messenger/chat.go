package messenger

import "time"

type Chat struct {
	Messages []Message `json:"messages"`
}

func NewChat() *Chat {
	return &Chat{
		Messages: []Message{},
	}
}

type Message struct {
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func (msg Message) isValid() bool {
	if len(msg.Username) > 0 && len(msg.Text) > 0 && len(msg.Text) < 300 {
		return true
	}
	return false
}

func (c *Chat) AddMessage(msg Message) {
	if msg.isValid() {
		c.Messages = append(c.Messages, msg)
	}
}
