package main

import (
	"github.com/ranefattesingh/ecommerce-platform/internal/config"
	"github.com/ranefattesingh/ecommerce-platform/internal/server/http"
	loader "github.com/ranefattesingh/ecommerce-platform/pkg/config"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// validate := validator.New(validator.WithRequiredStructEnabled())
	// h := handlers.NewUserHandler(validate)
	var configuration config.Config
	err := loader.LoadConfig("config.yaml", &configuration)
	if err != nil {
		logger.Fatal("config loading", zap.Error(err))
	}

	logger.Info("config load successful")

	httpServer := http.NewHTTPServer(configuration.Server)
	err = httpServer.StartServer()
	if err != nil {
		logger.Error("start server", zap.Error(err))
	}
}
