package queries

const GetChatById = `
SELECT
	c.id
	,c.name as chat_name
	,u.id as user_id
	,u.username as username
FROM chats c
JOIN chats_users cu ON c.id = cu.chat_id
LEFT JOIN messages m ON c.id = m.chat_id
JOIN users u ON cu.user_id = u.id
WHERE c.id = @chat_id
ORDER BY m.created_at DESC;`
