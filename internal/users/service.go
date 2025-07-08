package users

type IUsersService interface {
	GetAllUsers() ([]User, error)
	GetUserByID(id int64) (*User, error)
}

type UsersService struct {
	repository IUsersRepository
}

func NewUsersService(repository IUsersRepository) *UsersService {
	return &UsersService{repository: repository}
}

func (s *UsersService) GetAllUsers() ([]User, error) {
	users, err := s.repository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UsersService) GetUserByID(id int64) (*User, error) {
	user, err := s.repository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil // User not found
	}

	return user, nil
}
