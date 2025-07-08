package queries

const GetChatsForUser = `
SELECT u.id as receiver_id, u.username as receiver_username,
	m.content
FROM messages m
JOIN users u
	on (u.id = m.from_id)
WHERE
	m.to_id = @user_id
GROUP BY u.id;`
