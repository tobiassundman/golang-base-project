package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassundman/go-demo-app/internal/app/repository"
	"github.com/tobiassundman/go-demo-app/internal/app/service"
)

var (
	USER1_REPOSITORY = repository.User{
		ID:    1,
		Name:  "Name Name 1",
		Email: "email1@email.com",
		Age:   37,
	}
	USER2_REPOSITORY = repository.User{
		ID:    2,
		Name:  "Name Name 2",
		Email: "email2@email.com",
		Age:   102,
	}

	USER1_SERVICE = service.User{
		ID:    1,
		Name:  "Name Name 1",
		Email: "email1@email.com",
		Age:   37,
	}
	USER2_SERVICE = service.User{
		ID:    2,
		Name:  "Name Name 2",
		Email: "email2@email.com",
		Age:   102,
	}
)

var _ repository.UserRepository = &userRepositoryMock{}

type userRepositoryMock struct {
	GetAllFunc func() ([]*repository.User, error)
	GetFunc    func(id int) (*repository.User, error)
	CreateFunc func(user *repository.User) (int, error)
	UpdateFunc func(user *repository.User) error
	DeleteFunc func(id int) error
}

func (m *userRepositoryMock) GetAll() ([]*repository.User, error) {
	return m.GetAllFunc()
}

func (m *userRepositoryMock) Get(id int) (*repository.User, error) {
	return m.GetFunc(id)
}

func (m *userRepositoryMock) Create(user *repository.User) (int, error) {
	return m.CreateFunc(user)
}

func (m *userRepositoryMock) Update(user *repository.User) error {
	return m.UpdateFunc(user)
}

func (m *userRepositoryMock) Delete(id int) error {
	return m.DeleteFunc(id)
}

func TestGetAll(t *testing.T) {
	t.Parallel()
	t.Run("should return all users", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			GetAllFunc: func() ([]*repository.User, error) {
				return []*repository.User{&USER1_REPOSITORY, &USER2_REPOSITORY}, nil
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		users, err := userService.GetAll()
		require.NoError(t, err)

		// Assert
		assert.Len(t, users, 2)
		assert.Contains(t, users, &USER1_SERVICE)
		assert.Contains(t, users, &USER2_SERVICE)
	})

	t.Run("no users returns empty array", func(t *testing.T) {
		t.Parallel()
		// Arrange
		userRepositoryMock := &userRepositoryMock{
			GetAllFunc: func() ([]*repository.User, error) {
				return []*repository.User{}, nil
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		users, err := userService.GetAll()
		require.NoError(t, err)

		// Assert
		assert.Len(t, users, 0)
	})
}

func TestGet(t *testing.T) {
	t.Parallel()
	t.Run("should return user", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			GetFunc: func(id int) (*repository.User, error) {
				assert.Equal(t, 1, id)
				return &USER1_REPOSITORY, nil
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		user, err := userService.Get(1)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, &USER1_SERVICE, user)
	})

	t.Run("should return ErrUserNotFound", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			GetFunc: func(id int) (*repository.User, error) {
				return nil, repository.ErrUserNotFound
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		user, err := userService.Get(1)

		// Assert
		assert.Nil(t, user)
		assert.Equal(t, service.ErrUserNotFound, err)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()
	t.Run("should create user", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			CreateFunc: func(user *repository.User) (int, error) {
				assert.Equal(t, &USER1_REPOSITORY, user)
				return 1, nil
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		user, err := userService.Create(&USER1_SERVICE)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, &USER1_SERVICE, user)
	})

	t.Run("should return ErrUserAlreadyExists", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			CreateFunc: func(user *repository.User) (int, error) {
				return 0, repository.ErrUserAlreadyExists
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		_, err := userService.Create(&USER1_SERVICE)

		// Assert
		assert.Equal(t, service.ErrUserAlreadyExists, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	t.Run("should update user", func(t *testing.T) {
		t.Parallel()

		updateCalled := false

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			UpdateFunc: func(user *repository.User) error {
				assert.Equal(t, &USER1_REPOSITORY, user)
				updateCalled = true
				return nil
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		err := userService.Update(&USER1_SERVICE)
		require.NoError(t, err)

		// Assert
		assert.True(t, updateCalled)
	})

	t.Run("should return ErrUserNotFound", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			UpdateFunc: func(user *repository.User) error {
				return repository.ErrUserNotFound
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		err := userService.Update(&USER1_SERVICE)

		// Assert
		assert.Equal(t, service.ErrUserNotFound, err)
	})

	t.Run("should return ErrUserAlreadyExists", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			UpdateFunc: func(user *repository.User) error {
				return repository.ErrUserAlreadyExists
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		err := userService.Update(&USER1_SERVICE)

		// Assert
		assert.Equal(t, service.ErrUserAlreadyExists, err)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()
	t.Run("should delete user", func(t *testing.T) {
		t.Parallel()

		deleteCalled := false

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			DeleteFunc: func(id int) error {
				assert.Equal(t, 1, id)
				deleteCalled = true
				return nil
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		err := userService.Delete(1)
		require.NoError(t, err)

		// Assert
		assert.True(t, deleteCalled)
	})

	t.Run("should return ErrUserNotFound", func(t *testing.T) {
		t.Parallel()

		// Arrange
		userRepositoryMock := &userRepositoryMock{
			DeleteFunc: func(id int) error {
				return repository.ErrUserNotFound
			},
		}
		userService := service.NewUserService(userRepositoryMock)

		// Act
		err := userService.Delete(1)

		// Assert
		assert.Equal(t, service.ErrUserNotFound, err)
	})
}
