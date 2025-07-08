package chats

import "github.com/laps15/go-chat/internal/users"

type Message struct {
	ID        int64      `json:"id"`
	Sender    users.User `json:"sender" form:"sender"`
	Chat      Chat       `json:"chat" form:"chat"`
	Content   string     `json:"content" form:"content"`
	CreatedAt string     `json:"created_at" form:"created_at"`
	ReadAt    string     `json:"read_at" form:"read_at"`
}

type Chat struct {
	ID           int64                `json:"id"`
	Name         string               `json:"name"`
	Participants map[int64]users.User `json:"participants"`
	LastMessage  string               `json:"last_message"`
}
