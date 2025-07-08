package chats

import (
	"database/sql"

	"github.com/laps15/go-chat/internal/chats/internal/queries"
)

type IChatsRepository interface {
	CreateMessage(message *Message) (*Message, error)
	GetMessagesBySenderAndReceiver(senderID int64, receiverID int64) ([]Message, error)
	GetChatsForUser(userId int64) ([]Chat, error)
}

type ChatsRepository struct {
	db *sql.DB
}

func NewChatsRepository(db *sql.DB) *ChatsRepository {
	return &ChatsRepository{db: db}
}

func (mr *ChatsRepository) CreateMessage(message *Message) (*Message, error) {
	result, err := mr.db.Exec(queries.CreateMessageQuery,
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

func (mr *ChatsRepository) GetMessagesBySenderAndReceiver(senderID int64, receiverID int64) ([]Message, error) {
	rows, err := mr.db.Query(
		queries.GetMessagesOnChatQuery,
		sql.Named("first_user", senderID),
		sql.Named("second_user", receiverID),
		sql.Named("limit", 100),
		sql.Named("offset", 0))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(
			&message.ID,
			&message.Sender.ID,
			&message.Receiver.ID,
			&message.Content,
			&message.CreatedAt,
			&message.ReadAt,
			&message.Sender.Username,
			&message.Receiver.Username); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
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
		if err := rows.Scan(&chat.Receiver.ID, &chat.Receiver.Username, &chat.LastMessage); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
