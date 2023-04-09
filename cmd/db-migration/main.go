package main

import (
	"log"

	"github.com/tobiassundman/go-demo-app/pkg/environment"
	"github.com/tobiassundman/go-demo-app/pkg/logging"
	"github.com/tobiassundman/go-demo-app/pkg/test"
	"go.uber.org/zap"
)

var (
	dbUser     = environment.GetEnvOrDefault("DB_USER", "demo_user")
	dbPassword = environment.GetEnvOrDefault("DB_PASSWORD", "demo_password")
	dbHost     = environment.GetEnvOrDefault("DB_HOST", "localhost")
	dbPort     = environment.GetEnvOrDefault("DB_PORT", "5432")
	dbName     = environment.GetEnvOrDefault("DB_NAME", "demo_db")
)

func main() {
	logger, err := logging.NewProductionLogger()
	if err != nil {
		log.Fatal("Failed to create logger", err)
	}
	defer logger.Sync()

	migrationsDir := "db/migrations"
	err = test.Migrate(dbUser, dbPassword, dbHost, dbPort, dbName, migrationsDir)
	if err != nil {
		logger.Fatal("Failed to apply migrations", zap.Error(err))
	}
}
