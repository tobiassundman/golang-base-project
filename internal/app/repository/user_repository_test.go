package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassundman/go-demo-app/internal/app/repository"
	"github.com/tobiassundman/go-demo-app/pkg/test"
)

var (
	USER1 = repository.User{
		ID:    1,
		Name:  "Name Name 1",
		Email: "email1@email.com",
		Age:   37,
	}
	USER2 = repository.User{
		ID:    2,
		Name:  "Name Name 2",
		Email: "email2@email.com",
		Age:   102,
	}
)

func TestGetAll(t *testing.T) {
	t.Parallel()
	t.Run("should return all users", func(t *testing.T) {
		t.Parallel()

		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		_, err := pgRepository.Create(&USER1)
		require.NoError(t, err)
		_, err = pgRepository.Create(&USER2)
		require.NoError(t, err)

		// Act
		users, err := pgRepository.GetAll()
		require.NoError(t, err)

		// Assert
		assert.Len(t, users, 2)
		assert.Contains(t, users, &USER1)
		assert.Contains(t, users, &USER2)
	})

	t.Run("empty returns empty array", func(t *testing.T) {
		t.Parallel()
		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		// Act
		users, err := pgRepository.GetAll()
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
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		id, err := pgRepository.Create(&USER1)
		require.NoError(t, err)

		// Act
		user, err := pgRepository.Get(id)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, user, &USER1)
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()
		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		// Act
		_, err := pgRepository.Get(24)
		require.Error(t, err)

		// Assert
		assert.Equal(t, repository.ErrUserNotFound, err)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()
	t.Run("create user generates id", func(t *testing.T) {
		t.Parallel()

		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		user := repository.User{
			ID:    0,
			Name:  "Name Name 1",
			Email: "email1@email.com",
			Age:   37,
		}

		generatedID, err := pgRepository.Create(&user)
		require.NoError(t, err)

		// Act
		createdUser, err := pgRepository.Get(generatedID)
		require.NoError(t, err)

		// Assert
		assert.NotEqual(t, 0, generatedID)
		assert.Equal(t, createdUser.ID, generatedID)
		assert.Equal(t, createdUser.Name, user.Name)
		assert.Equal(t, createdUser.Email, user.Email)
		assert.Equal(t, createdUser.Age, user.Age)
	})

	t.Run("user already exists", func(t *testing.T) {
		t.Parallel()
		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		_, err := pgRepository.Create(&USER1)
		require.NoError(t, err)

		// Act
		_, err = pgRepository.Create(&USER1)
		require.Error(t, err)

		// Assert
		assert.Equal(t, repository.ErrUserAlreadyExists, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	t.Run("update existing", func(t *testing.T) {
		t.Parallel()

		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		id, err := pgRepository.Create(&USER1)
		require.NoError(t, err)

		modifiedUser := USER1
		modifiedUser.Name = "Modified Name"
		modifiedUser.Email = "Modified email"
		modifiedUser.Age = 99

		// Act
		err = pgRepository.Update(&modifiedUser)
		require.NoError(t, err)
		updatedUser, err := pgRepository.Get(id)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, &modifiedUser, updatedUser)
	})

	t.Run("update non-existing", func(t *testing.T) {
		t.Parallel()
		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		// Act
		err := pgRepository.Update(&USER1)
		require.Error(t, err)

		// Assert
		assert.Equal(t, repository.ErrUserNotFound, err)
	})

	t.Run("cannot use email of another user", func(t *testing.T) {
		t.Parallel()
		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		_, err := pgRepository.Create(&USER1)
		require.NoError(t, err)
		_, err = pgRepository.Create(&USER2)
		require.NoError(t, err)

		modifiedUser := USER2
		modifiedUser.Email = USER1.Email

		// Act
		err = pgRepository.Update(&modifiedUser)
		require.Error(t, err)

		// Assert
		assert.Equal(t, repository.ErrUserAlreadyExists, err)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()
	t.Run("delete existing", func(t *testing.T) {
		t.Parallel()

		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		id, err := pgRepository.Create(&USER1)
		require.NoError(t, err)

		// Act
		err = pgRepository.Delete(id)
		require.NoError(t, err)
		_, err = pgRepository.Get(id)
		require.Error(t, err)

		// Assert
		assert.Equal(t, err, repository.ErrUserNotFound)
	})

	t.Run("delete non-existing", func(t *testing.T) {
		t.Parallel()

		// Arrange
		db := test.StartDatabase(t)
		defer db.Close()
		pgRepository := repository.NewPostgresUserRepository(db, time.Second*2)

		// Act
		err := pgRepository.Delete(25)
		require.Error(t, err)

		// Assert
		assert.Equal(t, err, repository.ErrUserNotFound)
	})
}
