package chats

import (
	"fmt"

	"github.com/laps15/go-chat/internal/users"
)

type IChatsService interface {
	CreatePrivateChat(...int64) (*Chat, error)
	CreateChat(chatName string, userIds ...int64) (*Chat, error)
	SendMessage(from int64, to int64, content string) (*Message, error)
	GetMessagesForUser(users.User) ([]Message, error)
}

type ChatsService struct {
	repository IChatsRepository
}

func NewChatsService(repository IChatsRepository) *ChatsService {
	return &ChatsService{repository: repository}
}

func (ms *ChatsService) CreatePrivateChat(userIds ...int64) (*Chat, error) {
	if len(userIds) != 2 {
		return nil, fmt.Errorf("Two users are required to create private a chat")
	}
	chat := &Chat{
		Participants: []users.User{users.User{ID: userIds[0]}, users.User{ID: userIds[1]}},
	}

	return ms.repository.CreateChat(chat)
}

func (ms *ChatsService) CreateChat(chatName string, userIds ...int64) (*Chat, error) {
	if len(userIds) < 2 {
		return nil, fmt.Errorf("at least two users are required to create a chat")
	}
	chat := &Chat{
		Name:         chatName,
		Participants: make([]users.User, 0, len(userIds)),
	}

	for _, userId := range userIds {
		chat.Participants = append(chat.Participants, users.User{ID: userId})
	}

	return ms.repository.CreateChat(chat)
}

func (ms *ChatsService) SendMessage(from int64, to int64, content string) (*Message, error) {
	message := &Message{
		Sender:   users.User{ID: from},
		Receiver: users.User{ID: to},
		Content:  content,
	}

	return ms.repository.SendMessage(message)
}

func (ms *ChatsService) GetChatsForUser(user *users.User) ([]Chat, error) {
	return ms.repository.GetChatsForUser(user.ID)
}
