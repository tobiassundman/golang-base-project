package service

import (
	"errors"

	"github.com/tobiassundman/go-demo-app/internal/app/repository"
)

// UserService is the service for the user resource.
type UserService interface {
	// GetAll gets all users.
	GetAll() ([]*User, error)
	// Get gets a user by id.
	Get(id int) (*User, error)
	// Create creates a user.
	Create(user *User) (*User, error)
	// Update updates a user.
	Update(user *User) error
	// Delete deletes a user.
	Delete(id int) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		userRepository: repository,
	}
}

// GetAll gets all users.
func (s *userService) GetAll() ([]*User, error) {
	users, err := s.userRepository.GetAll()
	if err != nil {
		return nil, err
	}
	serviceUsers := make([]*User, len(users))
	for i, user := range users {
		serviceUsers[i] = repositoryUserToServiceUser(user)
	}
	return serviceUsers, nil
}

// Get gets a user by id.
func (s *userService) Get(id int) (*User, error) {
	user, err := s.userRepository.Get(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return repositoryUserToServiceUser(user), nil
}

// Create creates a user.
func (s *userService) Create(user *User) (*User, error) {
	id, err := s.userRepository.Create(serviceUserToRepositoryUser(user))
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	createdUser := &User{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
	return createdUser, nil
}

// Update updates a user.
func (s *userService) Update(user *User) error {
	err := s.userRepository.Update(serviceUserToRepositoryUser(user))
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(err, repository.ErrUserAlreadyExists):
		return ErrUserAlreadyExists
	}
	return err
}

// Delete deletes a user.
func (s *userService) Delete(id int) error {
	err := s.userRepository.Delete(id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return ErrUserNotFound
	}
	return err
}
