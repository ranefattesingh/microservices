package main

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/ranefattesingh/ecommerce-platform/internal/config"
	"github.com/ranefattesingh/ecommerce-platform/internal/platform/database/psql"
	"github.com/ranefattesingh/ecommerce-platform/internal/router"
	"github.com/ranefattesingh/ecommerce-platform/internal/server/http"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/handlers"
	userRepoPSQL "github.com/ranefattesingh/ecommerce-platform/internal/user/repository/psql"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/service"
	loader "github.com/ranefattesingh/ecommerce-platform/pkg/config"
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
	usersHandler := handlers.NewUserHandler(validate, service)

	v1Router := router.NewV1Router(usersHandler)

	logger.Info("config load successful")

	httpServer := http.NewHTTPServer(configuration.Server)
	err = httpServer.StartServer(v1Router)
	if err != nil {
		logger.Error("start server", zap.Error(err))
	}
}
