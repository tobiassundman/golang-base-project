package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

const (
	postgresGetAllUsersQuery = `SELECT id, name, email, age FROM config.users`
	postgresGetUserQuery     = `SELECT id, name, email, age FROM config.users WHERE id = $1`
	postgresCreateUserQuery  = `INSERT INTO config.users (name, email, age) VALUES ($1, $2, $3) RETURNING id`
	postgresUpdateUserQuery  = `UPDATE config.users SET name = $1, email = $2, age = $3 WHERE id = $4`
	postgresDeleteUserQuery  = `DELETE FROM config.users WHERE id = $1`
)

// UserRepository is an interface for the user repository
type UserRepository interface {
	// GetAll returns all users
	GetAll() ([]*User, error)
	// Get returns a user with the given id
	Get(id int) (*User, error)
	// Create creates a new user
	Create(user *User) (int, error)
	// Update updates a user
	Update(user *User) error
	// Delete deletes a user
	Delete(id int) error
}

// PostgresUserRepository is a repository for users in a Postgres database
type PostgresUserRepository struct {
	queryTimeout time.Duration
	db           *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB, queryTimeout time.Duration) *PostgresUserRepository {
	return &PostgresUserRepository{
		queryTimeout: queryTimeout,
		db:           db,
	}
}

// GetAll returns all users
func (r *PostgresUserRepository) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.queryTimeout)
	defer cancel()
	users := []*User{}
	err := r.db.SelectContext(ctx, &users, postgresGetAllUsersQuery)
	return users, err
}

// Get returns a user with the given id
func (r *PostgresUserRepository) Get(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.queryTimeout)
	defer cancel()
	user := &User{}
	err := r.db.GetContext(ctx, user, postgresGetUserQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	return user, err
}

// Create creates a new user
func (r *PostgresUserRepository) Create(user *User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.queryTimeout)
	defer cancel()
	var id int
	err := r.db.QueryRowContext(ctx, postgresCreateUserQuery, user.Name, user.Email, user.Age).Scan(&id)
	switch typedErr := err.(type) {
	case pgx.PgError:
		if typedErr.Code == pgerrcode.UniqueViolation {
			return 0, ErrUserAlreadyExists
		}
	}

	return id, err
}

// Update updates a user
func (r *PostgresUserRepository) Update(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.queryTimeout)
	defer cancel()
	result, err := r.db.ExecContext(ctx, postgresUpdateUserQuery, user.Name, user.Email, user.Age, user.ID)
	switch typedErr := err.(type) {
	case pgx.PgError:
		if typedErr.Code == pgerrcode.UniqueViolation {
			return ErrUserAlreadyExists
		}
	}
	if noRowsAffected(result) {
		return ErrUserNotFound
	}
	return nil
}

// Delete deletes a user
func (r *PostgresUserRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.queryTimeout)
	defer cancel()
	result, err := r.db.ExecContext(ctx, postgresDeleteUserQuery, id)
	if err != nil {
		return err
	}
	if noRowsAffected(result) {
		return ErrUserNotFound
	}
	return nil
}

// noRowsAffected returns true if the result of an update or delete query did not affect any rows (i.e. because the row did not exist)
func noRowsAffected(result sql.Result) bool {
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected == 0
}
