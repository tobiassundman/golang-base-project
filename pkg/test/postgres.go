package test

import (
	"database/sql"
	"testing"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"
	"github.com/tobiassundman/go-demo-app/pkg/database"
)

// StartDatabase starts a Postgres database in a Docker container returning a connection string.
func StartDatabase(t testing.TB) *sqlx.DB {

	env := []string{
		"POSTGRES_USER=demo_user",
		"POSTGRES_PASSWORD=demo_password",
		"POSTGRES_DB=demo_db",
	}

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env:        env,
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	require.NoError(t, err)

	resource.Expire(120)

	t.Cleanup(func() {
		err := pool.Purge(resource)
		require.NoError(t, err)
	})

	exposedPort := resource.GetPort("5432/tcp")

	db, err := database.UserDatabaseConnection("localhost", exposedPort, "demo_user", "demo_password", "demo_db")
	require.NoError(t, err)

	err = Migrate("demo_user", "demo_password", "localhost", exposedPort, "demo_db", "../../../db/migrations")
	require.NoError(t, err)

	return db
}

// Connect creates a connection to the database using the provided connection string.
func Connect(t testing.TB, connectString string) *sql.DB {
	connectConfig, err := pgx.ParseConnectionString(connectString)
	require.NoError(t, err)

	return stdlib.OpenDB(connectConfig)
}
