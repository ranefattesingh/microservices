package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/ranefattesingh/ecommerce-platform/auth/config"
	"github.com/ranefattesingh/ecommerce-platform/auth/platform/database/psql"
	"github.com/ranefattesingh/ecommerce-platform/auth/router"
	"github.com/ranefattesingh/ecommerce-platform/auth/server/http"
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

	_, err = psql.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		logger.Fatal("database connection", zap.Error(err))
	}

	v1Router := router.NewV1Router()

	logger.Info("config load successful")

	httpServer := http.NewHTTPServer(configuration.Server)
	err = httpServer.StartServer(v1Router)
	if err != nil {
		logger.Error("start server", zap.Error(err))
	}
}
