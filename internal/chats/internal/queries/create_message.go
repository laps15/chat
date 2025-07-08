package queries

const CreateMessageQuery = `
	INSERT INTO messages (chat_id, sender_id, content)
	VALUES (@chat_id, @sender_id, @content)`
