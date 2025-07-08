package chats

import (
	"fmt"

	"github.com/laps15/go-chat/internal/users"
)

type IChatsService interface {
	CreatePrivateChat(...int64) (*Chat, error)
	CreateChat(chatName string, userIds ...int64) (*Chat, error)
	SendMessage(from int64, to int64, content string) (*Message, error)
	GetPrivateChatForUsers(userIds ...int64) (*Chat, error)
	GetChatById(chatId int64) (*Chat, error)
	GetMessagesForChat(chat *Chat) []Message
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
		Participants: map[int64]users.User{
			userIds[0]: users.User{ID: userIds[0]},
			userIds[1]: users.User{ID: userIds[1]},
		},
	}

	return ms.repository.CreateChat(chat)
}

func (ms *ChatsService) CreateChat(chatName string, userIds ...int64) (*Chat, error) {
	if len(userIds) < 2 {
		return nil, fmt.Errorf("at least two users are required to create a chat")
	}
	chat := &Chat{
		Name:         chatName,
		Participants: make(map[int64]users.User),
	}

	for _, userId := range userIds {
		chat.Participants[userId] = users.User{ID: userId}
	}

	return ms.repository.CreateChat(chat)
}

func (ms *ChatsService) SendMessage(from int64, chat Chat, content string) (*Message, error) {
	message := &Message{
		Sender:  users.User{ID: from},
		Chat:    chat,
		Content: content,
	}

	return ms.repository.SendMessage(message)
}

func (ms *ChatsService) GetChatsForUser(user *users.User) ([]Chat, error) {
	return ms.repository.GetChatsForUser(user.ID)
}

func (ms *ChatsService) GetPrivateChatForUsers(userIds ...int64) (*Chat, error) {
	if len(userIds) != 2 {
		return nil, fmt.Errorf("Two users are required to get a private chat")
	}

	chat, err := ms.repository.GetChatForUsers(userIds[0], userIds[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get private chat for users %v and %v: %w", userIds[0], userIds[1], err)
	}

	return chat, nil
}

func (ms *ChatsService) GetChatById(chatId int64) (*Chat, error) {
	return ms.repository.GetChatById(chatId)
}

func (ms *ChatsService) GetMessagesForChat(chat *Chat) []Message {
	messages, err := ms.repository.GetMessagesForChat(chat)
	if err != nil {
		return []Message{}
	}

	return messages
}
