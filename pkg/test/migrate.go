package test

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrate runs the database migrations in the given directory during tests.
func Migrate(dbUser, dbPassword, dbHost, dbPort, dbName, migrationsDir string) error {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsDir),
		connectionString)
	if err != nil {
		return err
	}
	return m.Up()
}
