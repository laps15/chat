package chats

import (
	"context"
	"database/sql"

	"github.com/laps15/go-chat/internal/chats/internal/queries"
	"github.com/laps15/go-chat/internal/users"
)

type IChatsRepository interface {
	CreateChat(chat *Chat) (*Chat, error)
	SendMessage(message *Message) (*Message, error)
	GetChatsForUser(userId int64) ([]Chat, error)
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
	result, err := mr.db.Exec(queries.CreateChatQuery,
		sql.Named("from_id", message.Sender.ID),
		sql.Named("to_id", message.Receiver.ID),
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
		var me, receiver users.User
		if err := rows.Scan(&me.ID, &me.Username, &chat.Name, &receiver.ID, &receiver.Username, &chat.LastMessage); err != nil {
			return nil, err
		}
		chat.Participants = []users.User{receiver}

		if chat.Name == "" {
			chat.Name = receiver.Username
		}

		chats = append(chats, chat)
	}

	return chats, nil
}
