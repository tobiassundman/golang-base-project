package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tobiassundman/go-demo-app/internal/app/controller"
	"github.com/tobiassundman/go-demo-app/internal/app/repository"
	"github.com/tobiassundman/go-demo-app/internal/app/service"
	"github.com/tobiassundman/go-demo-app/pkg/database"
	"github.com/tobiassundman/go-demo-app/pkg/environment"
	"github.com/tobiassundman/go-demo-app/pkg/logging"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

var (
	serverPort   = environment.GetEnvOrDefault("SERVER_PORT", "8080")
	dbUser       = environment.GetEnvOrDefault("DB_USER", "demo_user")
	dbPassword   = environment.GetEnvOrDefault("DB_PASSWORD", "demo_password")
	dbHost       = environment.GetEnvOrDefault("DB_HOST", "localhost")
	dbPort       = environment.GetEnvOrDefault("DB_PORT", "5432")
	dbName       = environment.GetEnvOrDefault("DB_NAME", "demo_db")
	queryTimeout = environment.GetEnvOrDefault("QUERY_TIMEOUT", "5s")
)

func main() {
	logger, err := logging.NewProductionLogger()
	if err != nil {
		log.Fatal("Failed to create logger", err)
	}
	defer logger.Sync()

	logger.Info("Starting demo app", zap.String("port", serverPort), zap.String("dbHost", dbHost), zap.String("dbPort", dbPort))
	db, err := database.UserDatabaseConnection(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	parsedQueryTimeout, err := time.ParseDuration(queryTimeout)
	if err != nil {
		logger.Fatal("Failed to parse query timeout", zap.Error(err))
	}

	userRepository := repository.NewPostgresUserRepository(db, parsedQueryTimeout)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService, logger)

	router := createRouter(logger)
	userController.ConfigureRoutes(router)

	p := ginprometheus.NewPrometheus("gin")

	p.Use(router)

	router.GET("/liveness", liveness)
	router.GET("/readiness", readiness(db))

	runServer(router, logger.Sugar())
}

// createRouter creates a new gin router with middleware
func createRouter(logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	return router
}

// runServer starts the http server and handles graceful shutdown
func runServer(router *gin.Engine, logger *zap.SugaredLogger) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", serverPort),
		Handler: router,
	}

	sigtermChannel := make(chan os.Signal, 1)
	signal.Notify(sigtermChannel, syscall.SIGTERM)
	defer close(sigtermChannel)

	shutdownWaitGroup := sync.WaitGroup{}
	shutdownWaitGroup.Add(1)
	go func() {
		defer shutdownWaitGroup.Done()
		<-sigtermChannel
		logger.Info("SIGTERM received, shutting down http server")

		shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			if err == http.ErrServerClosed {
				logger.Info("Graceful shutdown of http server initiated")
				return
			}
			logger.Error("Failed to shut down http server", err)
			os.Exit(1)
		}
	}()

	shutdownWaitGroup.Add(1)
	go func() {
		defer shutdownWaitGroup.Done()

		logger.Info("Starting http server")
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal("Failed to start http server", err)
		}
	}()

	shutdownWaitGroup.Wait()
}

// liveness checks if the service has started
func liveness(c *gin.Context) {
	c.Status(http.StatusOK)
}

// readiness checks if the application is ready to accept requests
func readiness(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		if db.Ping() != nil {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		c.Status(http.StatusOK)
	}
}
