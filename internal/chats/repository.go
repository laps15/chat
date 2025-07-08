package chats

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/laps15/go-chat/internal/chats/internal/queries"
	"github.com/laps15/go-chat/internal/users"
)

type IChatsRepository interface {
	CreateChat(chat *Chat) (*Chat, error)
	SendMessage(message *Message) (*Message, error)
	GetChatsForUser(userId int64) ([]Chat, error)
	GetChatForUsers(userIds ...int64) (*Chat, error)
	GetChatById(userId int64, chatId int64) (*Chat, error)
	GetMessagesForChat(chat *Chat) ([]Message, error)
}

type ChatsRepository struct {
	db *sql.DB
}

func NewChatsRepository(db *sql.DB) *ChatsRepository {
	return &ChatsRepository{db: db}
}

func (mr *ChatsRepository) CreateChat(chat *Chat) (*Chat, error) {
	tx, err := mr.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	res, err := mr.db.Exec(queries.CreateChatQuery,
		sql.Named("chat_name", chat.Name))
	if err != nil {
		return nil, err
	}

	chatId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	for _, participant := range chat.Participants {
		_, err := mr.db.Exec(queries.AddUserToChatQuery,
			sql.Named("chat_id", chatId),
			sql.Named("user_id", participant.ID))
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	return chat, nil
}

func (mr *ChatsRepository) SendMessage(message *Message) (*Message, error) {
	result, err := mr.db.Exec(queries.CreateMessageQuery,
		sql.Named("chat_id", message.Chat.ID),
		sql.Named("sender_id", message.Sender.ID),
		sql.Named("content", message.Content))
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	message.ID = id

	return message, nil
}

func (mr *ChatsRepository) GetChatsForUser(userId int64) ([]Chat, error) {
	rows, err := mr.db.Query(
		queries.GetChatsForUser,
		sql.Named("user_id", userId))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		chat.Participants = make(map[int64]users.User)
		var me, receiver users.User
		if err := rows.Scan(&chat.ID, &me.ID, &me.Username, &chat.Name, &receiver.ID, &receiver.Username, &chat.LastMessage); err != nil {
			return nil, err
		}
		chat.Participants[me.ID] = me
		chat.Participants[receiver.ID] = receiver

		if chat.Name == "" {
			chat.Name = receiver.Username
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (mr *ChatsRepository) GetChatForUsers(userIds ...int64) (*Chat, error) {
	strs := make([]string, len(userIds))
	for i, id := range userIds {
		strs[i] = strconv.FormatInt(id, 10)
	}

	query := fmt.Sprintf(queries.GetChatForUsersQuery, strings.Join(strs, ","))
	rows, err := mr.db.Query(
		query,
		sql.Named("user_ids_count", len(userIds)))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rowCount = 0
	var chat Chat
	chat.Participants = make(map[int64]users.User)
	for rows.Next() {
		rowCount++
		var user users.User
		if err := rows.Scan(&chat.ID, &chat.Name, &user.ID, &user.Username); err != nil {
			return nil, err
		}
		chat.Participants[user.ID] = user

		if chat.Name == "" {
			chat.Name = user.Username
		}
	}

	if rowCount == 0 {
		return nil, nil
	}

	if rowCount != len(userIds) {
		return nil, fmt.Errorf("expected %d users, but found %d in chat", len(userIds), rowCount)
	}

	return &chat, nil
}

func (mr *ChatsRepository) GetChatById(userId int64, chatId int64) (*Chat, error) {
	rows, err := mr.db.Query(
		queries.GetChatById,
		sql.Named("chat_id", chatId))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernamesInChat []string
	var chat Chat
	chat.Participants = make(map[int64]users.User)
	for rows.Next() {
		var user users.User
		if err := rows.Scan(&chat.ID, &chat.Name, &user.ID, &user.Username); err != nil {
			return nil, err
		}
		chat.Participants[user.ID] = user

		if user.ID != userId {
			usernamesInChat = append(usernamesInChat, user.Username)
		}
	}

	if chat.Name == "" && len(usernamesInChat) > 0 {
		chat.Name = strings.Join(usernamesInChat, ", ")
	}

	return &chat, nil
}

func (mr *ChatsRepository) GetMessagesForChat(chat *Chat) ([]Message, error) {
	rows, err := mr.db.Query(
		queries.GetMessagesForChatQuery,
		sql.Named("chat_id", chat.ID),
		sql.Named("limit", 100),
		sql.Named("offset", 0))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.Content, &message.CreatedAt, &message.Sender.ID, &message.Sender.Username); err != nil {
			return nil, err
		}
		message.Chat = *chat
		messages = append(messages, message)
	}

	return messages, nil
}
