package auth

import (
	"errors"
	"strings"

	"github.com/labstack/gommon/random"
	"github.com/laps15/go-chat/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	CreateUser(user *users.User) (*users.User, error)
	AuthenticateUser(username string, password string) (*users.User, error)
	CreateSession(userID int64) (string, error)
}

type AuthService struct {
	repository users.IUsersRepository
}

func NewAuthService(repository users.IUsersRepository) *AuthService {
	return &AuthService{repository: repository}
}

func getPasswordSalt() string {
	return random.String(16)
}

func getHashedPassword(password string, salt string) (string, error) {
	hpwd, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hpwd) + ":" + salt, nil
}

func (s *AuthService) CreateUser(user *users.User) (*users.User, error) {
	hpwd, err := getHashedPassword(user.Password, getPasswordSalt())
	if err != nil {
		return nil, err
	}

	user.Password = hpwd

	return s.repository.CreateUser(user)
}

func (s *AuthService) AuthenticateUser(username string, password string) (*users.User, error) {
	user, err := s.repository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("UserNotFound")
	}

	hpwd, salt := user.Password[:strings.LastIndex(user.Password, ":")], user.Password[strings.LastIndex(user.Password, ":")+1:]
	if err := bcrypt.CompareHashAndPassword([]byte(hpwd), []byte(password+salt)); err != nil {
		return nil, errors.New("InvalidCredentials")
	}

	return user, nil
}
