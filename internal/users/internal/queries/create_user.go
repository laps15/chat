package queries

const CreateUserQuery = `
INSERT INTO users (username, password, email)
VALUES (?, ?, ?);`
