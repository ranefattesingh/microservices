package main

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	loader "github.com/ranefattesingh/ecommerce-platform/pkg/config"
	"github.com/ranefattesingh/ecommerce-platform/users/config"
	"github.com/ranefattesingh/ecommerce-platform/users/handlers"
	"github.com/ranefattesingh/ecommerce-platform/users/platform/database/psql"
	userRepoPSQL "github.com/ranefattesingh/ecommerce-platform/users/repository/psql"
	"github.com/ranefattesingh/ecommerce-platform/users/router"
	"github.com/ranefattesingh/ecommerce-platform/users/server/http"
	"github.com/ranefattesingh/ecommerce-platform/users/service"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("env config loading", zap.Error(err))
	}

	var configuration config.Config
	err = loader.LoadConfig("config.yaml", &configuration)
	if err != nil {
		logger.Fatal("yaml config loading", zap.Error(err))
	}

	db, err := psql.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		logger.Fatal("database connection", zap.Error(err))
	}

	repo := userRepoPSQL.NewUsersRepository(db)
	service := service.NewUsersService(repo)
	validate := validator.New(validator.WithRequiredStructEnabled())
	usersHandler := handlers.NewUserHandler(validate, service, logger)

	v1Router := router.NewV1Router(usersHandler)

	logger.Info("config load successful")

	httpServer := http.NewHTTPServer(configuration.Server)
	err = httpServer.StartServer(v1Router)
	if err != nil {
		logger.Error("start server", zap.Error(err))
	}
}
