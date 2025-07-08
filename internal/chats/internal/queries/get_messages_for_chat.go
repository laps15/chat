package queries

const GetMessagesForChatQuery = `
	SELECT m.id
	,m.content
	,m.created_at
	,m.sender_id
	,u.username AS sender_username
FROM chats c
JOIN messages m ON c.id = m.chat_id
JOIN users u ON m.sender_id = u.id
WHERE c.id = @chat_id
ORDER BY m.created_at
LIMIT @limit OFFSET @offset;`
