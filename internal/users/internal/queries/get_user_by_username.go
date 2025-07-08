package queries

const GetUserByUsernameQuery = `
SELECT id, username, password, email
FROM users
WHERE username = ?;`
