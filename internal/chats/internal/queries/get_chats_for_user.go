package queries

const GetChatsForUser = `
SELECT 
	me.id as my_id
	,me.username as my_username
    ,c.name as chat_name
    ,other.id as receiver_id
    ,other.username as receiver_username
    ,COALESCE(m.content, "") as last_message
FROM chats c
JOIN chats_users cu ON c.id = cu.chat_id
LEFT JOIN messages m ON c.id = m.chat_id
JOIN users me ON cu.user_id = me.id
JOIN chats_users other_cu ON c.id = other_cu.chat_id AND other_cu.user_id != me.id
JOIN users other ON other_cu.user_id = other.id
WHERE me.id = @user_id
GROUP BY c.id, me.id, other.id
ORDER BY m.created_at DESC;`
