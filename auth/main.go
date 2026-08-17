package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ranefattesingh/microservices/auth/config"
	"github.com/ranefattesingh/microservices/auth/handler"
	"github.com/ranefattesingh/microservices/auth/server/http"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.DefaultLoader().Load()
	if err != nil {
		logger.Fatal("fail to load config", zap.Error(err))
	}

	authHandler := handler.NewAuthHandler(logger)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	server := http.NewHTTPServer(logger, cfg.Server)

	go func() {
		if err := server.Start(authHandler); err != nil {
			logger.Error("HTTP server failed", zap.Error(err))
			stop()
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(); err != nil {
		logger.Error("HTTP server shutdown failed", zap.Error(err))
	}
}
