package queries

const CreateMessageQuery = `
INSERT INTO messages (from_id, to_id, content)
VALUES (@from_id, @to_id, @content);`
