package queries

const GetMessagesOnChatQuery = `
SELECT m.id, m.from_id, m.to_id, m.content, m.created_at,
	   m.read_at, u1.username AS sender_username, 
	   u2.username AS receiver_username
FROM messages m
JOIN users u1 ON m.from_id = u1.id
JOIN users u2 ON m.to_id = u2.id
WHERE (m.from_id = @first_user AND m.to_id = @second_user) OR (m.from_id = @second_user AND m.to_id = @first_user)
ORDER BY m.created_at DESC
LIMIT @limit OFFSET @offset;`
