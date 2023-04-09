package controller

import "github.com/tobiassundman/go-demo-app/internal/app/service"

// User is the user model for the controller layer.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
}

// UpdateUserRequest is the request model for updating a user.
type UpdateUserRequest struct {
	ID    int    `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
}

// GetUsersResponse is the response model when getting all users.
type GetUsersResponse struct {
	Users []*User `json:"users"`
}

// updateUserRequestToServiceUser converts a controller UpdateUserRequest to a service User.
func updateUserRequestToServiceUser(user *UpdateUserRequest) *service.User {
	return &service.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}

// createUserRequestToServiceUser converts a controller User to a service User.
func createUserRequestToServiceUser(user *User) *service.User {
	return &service.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}

// serviceUserToControllerUser converts a service User to a controller User.
func serviceUserToControllerUser(user *service.User) *User {
	return &User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   user.Age,
	}
}
