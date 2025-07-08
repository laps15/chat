package queries

const GetUserByIDQuery = `
SELECT id, username, email
FROM users
WHERE id = @user_id;`
