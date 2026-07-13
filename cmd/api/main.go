package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/ecommerce-platform/internal/config"
	"github.com/ranefattesingh/ecommerce-platform/internal/router"
	"github.com/ranefattesingh/ecommerce-platform/internal/server/http"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/handlers"
	loader "github.com/ranefattesingh/ecommerce-platform/pkg/config"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	validate := validator.New(validator.WithRequiredStructEnabled())
	usersHandler := handlers.NewUserHandler(validate)

	v1Router := router.NewV1Router(usersHandler)

	var configuration config.Config
	err := loader.LoadConfig("config.yaml", &configuration)
	if err != nil {
		logger.Fatal("config loading", zap.Error(err))
	}

	logger.Info("config load successful")

	httpServer := http.NewHTTPServer(configuration.Server)
	err = httpServer.StartServer(v1Router)
	if err != nil {
		logger.Error("start server", zap.Error(err))
	}
}
