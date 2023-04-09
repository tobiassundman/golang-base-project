package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tobiassundman/go-demo-app/internal/app/service"
	"go.uber.org/zap"
)

// UserCtrl is the controller for the user resource.
type UserController struct {
	logger      *zap.Logger
	userService service.UserService
}

func NewUserController(service service.UserService, logger *zap.Logger) *UserController {
	return &UserController{
		logger:      logger,
		userService: service,
	}
}

// ConfigureRoutes configures the routes for the user resource.
func (c *UserController) ConfigureRoutes(router *gin.Engine) {
	userGroup := router.Group("/v1")
	userGroup.GET("/users", c.getUsers)
	userGroup.GET("/users/:id", c.getUser)
	userGroup.POST("/users", c.createUser)
	userGroup.PUT("/users", c.updateUser)
	userGroup.DELETE("/users/:id", c.deleteUser)
}

// getUsers returns all users.
func (c *UserController) getUsers(ctx *gin.Context) {
	users, err := c.userService.GetAll()
	if err != nil {
		c.logger.Error("Failed to get users", zap.Error(err))
		apiError := apiErrorFromServiceError(err)
		ctx.JSON(apiError.Status, apiError)
		return
	}

	usersInResponse := make([]*User, len(users))
	for i, user := range users {
		usersInResponse[i] = serviceUserToControllerUser(user)
	}

	reponse := GetUsersResponse{
		Users: usersInResponse,
	}

	ctx.JSON(http.StatusOK, reponse)
}

// getUser returns a single user by id.
func (c *UserController) getUser(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		c.logger.Warn("Failed to parse id", zap.Error(err))
		ctx.JSON(ErrInvalidID.Status, ErrInvalidID)
		return
	}

	user, err := c.userService.Get(parsedID)
	if err != nil {
		apiError := apiErrorFromServiceError(err)
		if apiError != ErrUserNotFound {
			c.logger.Warn("Failed to get user", zap.Error(err), zap.Int("id", parsedID))
		}
		ctx.JSON(apiError.Status, apiError)
		return
	}

	ctx.JSON(http.StatusOK, serviceUserToControllerUser(user))
}

// createUser creates a new user.
func (c *UserController) createUser(ctx *gin.Context) {
	inputUser := User{}
	err := ctx.BindJSON(&inputUser)
	if err != nil {
		c.logger.Warn("Failed to parse user", zap.Error(err))
		ctx.JSON(ErrValidationFailed.Status, ErrValidationFailed)
		return
	}
	newUser, err := c.userService.Create(createUserRequestToServiceUser(&inputUser))
	if err != nil {
		c.logger.Warn("Failed to create user", zap.Error(err), zap.Any("user", inputUser))
		apiError := apiErrorFromServiceError(err)
		ctx.JSON(apiError.Status, apiError)
		return
	}

	ctx.JSON(http.StatusCreated, serviceUserToControllerUser(newUser))
}

// updateUser updates an existing user by id.
func (c *UserController) updateUser(ctx *gin.Context) {
	inputUser := UpdateUserRequest{}
	err := ctx.BindJSON(&inputUser)
	if err != nil {
		c.logger.Warn("Failed to parse user", zap.Error(err))
		ctx.JSON(ErrValidationFailed.Status, ErrValidationFailed)
		return
	}

	err = c.userService.Update(updateUserRequestToServiceUser(&inputUser))
	if err != nil {
		apiError := apiErrorFromServiceError(err)
		if apiError != ErrUserNotFound {
			c.logger.Warn("Failed to update user", zap.Error(err), zap.Any("user", inputUser))
		}
		ctx.JSON(apiError.Status, apiError)
		return
	}

	ctx.Status(http.StatusOK)
}

// deleteUser deletes an existing user by id.
func (c *UserController) deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		c.logger.Warn("Failed to parse id", zap.Error(err), zap.String("id", id))
		ctx.JSON(ErrInvalidID.Status, ErrInvalidID)
		return
	}

	err = c.userService.Delete(parsedID)
	if err != nil {
		apiError := apiErrorFromServiceError(err)
		if apiError != ErrUserNotFound {
			c.logger.Warn("Failed to delete user", zap.Error(err), zap.Int("id", parsedID))
		}
		ctx.JSON(apiError.Status, apiError)
		return
	}

	ctx.Status(http.StatusOK)
}
