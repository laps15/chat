package queries

const GetChatForUsersQuery = `
SELECT c.id, c.name, u.id, u.username
FROM chats c
JOIN chats_users cu ON c.id = cu.chat_id
JOIN users u ON cu.user_id = u.id
WHERE c.id = (
	SELECT cu.chat_id
	FROM chats_users cu
	WHERE cu.user_id IN (%s)
	GROUP BY cu.chat_id
	HAVING COUNT(DISTINCT cu.user_id) = @user_ids_count
);`
