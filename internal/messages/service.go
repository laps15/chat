package messages

import "github.com/laps15/go-chat/internal/users"

type IMessagesService interface {
	CreateMessage(from int64, to int64, content string) (*Message, error)
	GetMessagesForUser(users.User) ([]Message, error)
}

type MessagesService struct {
	repository IMessagesRepository
}

func NewMessagesService(repository IMessagesRepository) *MessagesService {
	return &MessagesService{repository: repository}
}

func (ms *MessagesService) CreateMessage(from int64, to int64, content string) (*Message, error) {
	message := &Message{
		Sender:   users.User{ID: from},
		Receiver: users.User{ID: to},
		Content:  content,
	}

	return ms.repository.CreateMessage(message)
}

func (ms *MessagesService) GetChatsForUser(user *users.User) ([]Chat, error) {
	return ms.repository.GetChatsForUser(user.ID)
}
