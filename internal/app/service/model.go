package service

import "github.com/tobiassundman/go-demo-app/internal/app/repository"

// User is the user model for the service layer.
type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

// repositoryUserToServiceUser converts a repository User to a service User.
func repositoryUserToServiceUser(user *repository.User) *User {
	return &User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}

// serviceUserToRepositoryUser converts a service User to a repository User.
func serviceUserToRepositoryUser(user *User) *repository.User {
	return &repository.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}
