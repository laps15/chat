package chats

import "github.com/laps15/go-chat/internal/users"

type IChatsService interface {
	CreateMessage(from int64, to int64, content string) (*Message, error)
	GetMessagesForUser(users.User) ([]Message, error)
}

type ChatsService struct {
	repository IChatsRepository
}

func NewChatsService(repository IChatsRepository) *ChatsService {
	return &ChatsService{repository: repository}
}

func (ms *ChatsService) CreateMessage(from int64, to int64, content string) (*Message, error) {
	message := &Message{
		Sender:   users.User{ID: from},
		Receiver: users.User{ID: to},
		Content:  content,
	}

	return ms.repository.CreateMessage(message)
}

func (ms *ChatsService) GetChatsForUser(user *users.User) ([]Chat, error) {
	return ms.repository.GetChatsForUser(user.ID)
}
