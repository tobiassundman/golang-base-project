package controller_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassundman/go-demo-app/internal/app/controller"
	"github.com/tobiassundman/go-demo-app/internal/app/service"
	"go.uber.org/zap"
)

var _ service.UserService = &userServiceMock{}

type userServiceMock struct {
	GetAllFunc func() ([]*service.User, error)
	GetFunc    func(id int) (*service.User, error)
	CreateFunc func(user *service.User) (*service.User, error)
	UpdateFunc func(user *service.User) error
	DeleteFunc func(id int) error
}

func (m *userServiceMock) GetAll() ([]*service.User, error) {
	return m.GetAllFunc()
}

func (m *userServiceMock) Get(id int) (*service.User, error) {
	return m.GetFunc(id)
}

func (m *userServiceMock) Create(user *service.User) (*service.User, error) {
	return m.CreateFunc(user)
}

func (m *userServiceMock) Update(user *service.User) error {
	return m.UpdateFunc(user)
}

func (m *userServiceMock) Delete(id int) error {
	return m.DeleteFunc(id)
}

func TestGetAll(t *testing.T) {
	t.Run("returns all users", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			GetAllFunc: func() ([]*service.User, error) {
				return []*service.User{
					{
						ID:    1,
						Name:  "Name Name 1",
						Email: "email1@email.com",
						Age:   37,
					},
					{
						ID:    2,
						Name:  "Name Name 2",
						Email: "email2@email.com",
						Age:   102,
					},
				}, nil
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)

		r := gofight.New()

		// Act
		r.GET("/v1/users").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, r.Code)

				assert.JSONEq(t,
					`{
						"users": [
							{
								"id": 1,
								"name": "Name Name 1",
								"email": "email1@email.com",
								"age": 37
							},
							{
								"id": 2,
								"name": "Name Name 2",
								"email": "email2@email.com",
								"age": 102
							}
						]	
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns empty list when no users", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			GetAllFunc: func() ([]*service.User, error) {
				return []*service.User{}, nil
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)

		r := gofight.New()

		// Act
		r.GET("/v1/users").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, r.Code)

				assert.JSONEq(t,
					`{
						"users": []
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 500 when error", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			GetAllFunc: func() ([]*service.User, error) {
				return nil, errors.New("error")
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)

		r := gofight.New()

		// Act
		r.GET("/v1/users").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
				assert.JSONEq(t,
					`{
						"error_code": "ErrInternalServer", 
						"error_message": "internal server error",
						"status": 500
					}`,
					r.Body.String(),
				)
			})
	})
}

func TestGet(t *testing.T) {
	t.Run("returns user", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			GetFunc: func(id int) (*service.User, error) {
				assert.Equal(t, 1, id)
				return &service.User{
					ID:    1,
					Name:  "Name Name 1",
					Email: "email1@email.com",
					Age:   37,
				}, nil
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.GET("/v1/users/1").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, r.Code)

				assert.JSONEq(t,
					`{
						"id": 1,
						"name": "Name Name 1",
						"email": "email1@email.com",
						"age": 37
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 404 when user not found", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			GetFunc: func(id int) (*service.User, error) {
				assert.Equal(t, 1, id)
				return nil, service.ErrUserNotFound
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.GET("/v1/users/1").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusNotFound, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrUserNotFound", 
						"error_message": "user not found", 
						"status": 404
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 500 when error", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			GetFunc: func(id int) (*service.User, error) {
				return nil, errors.New("error")
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.GET("/v1/users/1").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrInternalServer",
						"error_message": "internal server error",
						"status": 500
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 400 when invalid id", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.GET("/v1/users/invalid").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusBadRequest, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrInvalidID",
						"error_message": "invalid id",
						"status": 400
					}`,
					r.Body.String(),
				)
			})
	})
}

func TestCreate(t *testing.T) {
	t.Run("creates user", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			CreateFunc: func(user *service.User) (*service.User, error) {
				return &service.User{
					ID:    1,
					Name:  "Name Name 1",
					Email: "email1@email.com",
					Age:   37,
				}, nil
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.POST("/v1/users").
			SetJSON(gofight.D{
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusCreated, r.Code)

				assert.JSONEq(t,
					`{
						"id": 1,
						"name": "Name Name 1",
						"email": "email1@email.com",
						"age": 37
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns conflict when user already exists", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			CreateFunc: func(user *service.User) (*service.User, error) {
				return nil, service.ErrUserAlreadyExists
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.POST("/v1/users").
			SetJSON(gofight.D{
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusConflict, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrUserAlreadyExists", 
						"error_message": "user already exists", 
						"status": 409
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 500 when error", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			CreateFunc: func(user *service.User) (*service.User, error) {
				return nil, errors.New("error")
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.POST("/v1/users").
			SetJSON(gofight.D{
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrInternalServer",
						"error_message": "internal server error",
						"status": 500
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 400 when invalid email", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.POST("/v1/users").
			SetJSON(gofight.D{
				"name":  "Name Name 1",
				"email": "invalid",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusBadRequest, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrValidationFailed",
						"error_message": "validation failed",
						"status": 400
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("additional fields are ignored", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			CreateFunc: func(user *service.User) (*service.User, error) {
				return &service.User{
					ID:    1,
					Name:  "Name Name 1",
					Email: "email1@email.com",
					Age:   37,
				}, nil
			},
		}

		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.POST("/v1/users").
			SetJSON(gofight.D{
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
				"foo":   "bar",
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusCreated, r.Code)

				assert.JSONEq(t,
					`{
						"id": 1,
						"name": "Name Name 1",
						"email": "email1@email.com",
						"age": 37
					}`,
					r.Body.String(),
				)
			})
	})

}

func TestUpdate(t *testing.T) {
	t.Run("updates user", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			UpdateFunc: func(user *service.User) error {
				assert.Equal(t, 1, user.ID)
				return nil
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.PUT("/v1/users").
			SetJSON(gofight.D{
				"id":    1,
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, r.Code)
			})
	})

	t.Run("returns 404 when user not found", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			UpdateFunc: func(user *service.User) error {
				return service.ErrUserNotFound
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.PUT("/v1/users").
			SetJSON(gofight.D{
				"id":    1,
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusNotFound, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrUserNotFound", 
						"error_message": "user not found", 
						"status": 404
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 500 when error", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			UpdateFunc: func(user *service.User) error {
				return errors.New("error")
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.PUT("/v1/users").
			SetJSON(gofight.D{
				"id":    1,
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrInternalServer",
						"error_message": "internal server error",
						"status": 500
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 400 when invalid email", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.PUT("/v1/users").
			SetJSON(gofight.D{
				"id":    1,
				"name":  "Name Name 1",
				"email": "invalid",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusBadRequest, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrValidationFailed",
						"error_message": "validation failed",
						"status": 400
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("additional fields are ignored", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			UpdateFunc: func(user *service.User) error {
				return nil
			},
		}

		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.PUT("/v1/users").
			SetJSON(gofight.D{
				"id":    1,
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
				"foo":   "bar",
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, r.Code)
			})
	})

	t.Run("conflict with other user", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			UpdateFunc: func(user *service.User) error {
				return service.ErrUserAlreadyExists
			},
		}

		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.PUT("/v1/users").
			SetJSON(gofight.D{
				"id":    1,
				"name":  "Name Name 1",
				"email": "email1@email.com",
				"age":   37,
			}).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusConflict, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrUserAlreadyExists", 
						"error_message": "user already exists", 
						"status": 409
					}`,
					r.Body.String(),
				)
			})
	})
}

func TestDelete(t *testing.T) {
	t.Run("deletes user", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			DeleteFunc: func(id int) error {
				assert.Equal(t, 1, id)
				return nil
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.DELETE("/v1/users/1").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, r.Code)
			})
	})

	t.Run("returns 404 when user not found", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			DeleteFunc: func(id int) error {
				return service.ErrUserNotFound
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.DELETE("/v1/users/1").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusNotFound, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrUserNotFound", 
						"error_message": "user not found", 
						"status": 404
					}`,
					r.Body.String(),
				)
			})
	})

	t.Run("returns 500 when error", func(t *testing.T) {
		t.Parallel()
		// Arrange
		serviceMock := &userServiceMock{
			DeleteFunc: func(id int) error {
				return errors.New("error")
			},
		}
		controller := controller.NewUserController(serviceMock, zap.NewNop())

		router := gin.Default()
		controller.ConfigureRoutes(router)
		r := gofight.New()

		// Act
		r.DELETE("/v1/users/1").
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
				require.JSONEq(
					t,
					`{
						"error_code": "ErrInternalServer",
						"error_message": "internal server error",
						"status": 500
					}`,
					r.Body.String(),
				)
			})
	})
}
