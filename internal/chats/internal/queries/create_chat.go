package queries

const CreateChatQuery = `
	INSERT INTO chats (name)
	VALUES (@chat_name)
`

const AddUserToChatQuery = `
INSERT INTO chats_users (chat_id, user_id)
VALUES (@chat_id, @user_id);`
