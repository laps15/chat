package users

import (
	"database/sql"

	"github.com/laps15/go-chat/internal/users/internal/queries"
)

type IUsersRepository interface {
	CreateUser(user *User) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetAllUsers() ([]User, error)
	GetUserByID(id int64) (*User, error)
}

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (ur *UsersRepository) CreateUser(user *User) (*User, error) {
	nuser, err := ur.db.Exec(queries.CreateUserQuery, user.Username, user.Password, user.Email)
	if err != nil {
		return nil, err
	}

	id, err := nuser.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.ID = id

	return user, nil
}

func (ur *UsersRepository) GetUserByUsername(username string) (*User, error) {
	row := ur.db.QueryRow(queries.GetUserByUsernameQuery, username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, err // Other error
	}

	return user, nil
}

func (ur *UsersRepository) GetAllUsers() ([]User, error) {
	rows, err := ur.db.Query(queries.GetAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UsersRepository) GetUserByID(id int64) (*User, error) {
	row := ur.db.QueryRow(queries.GetUserByIDQuery, sql.Named("user_id", id))

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
