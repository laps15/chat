package chats

import "github.com/laps15/go-chat/internal/users"

type Message struct {
	ID        int64      `json:"id"`
	Sender    users.User `json:"sender" form:"sender"`
	Receiver  users.User `json:"receiver" form:"receiver"`
	Content   string     `json:"content" form:"content"`
	CreatedAt string     `json:"created_at" form:"created_at"`
	ReadAt    string     `json:"read_at" form:"read_at"`
}

type Chat struct {
	Name         string       `json:"name"`
	Participants []users.User `json:"participants"`
	LastMessage  string       `json:"last_message"`
}
