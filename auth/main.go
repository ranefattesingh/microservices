package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ranefattesingh/microservices/auth/config"
	"github.com/ranefattesingh/microservices/auth/handler"
	"github.com/ranefattesingh/microservices/auth/repository"
	"github.com/ranefattesingh/microservices/auth/server/http"
	"github.com/ranefattesingh/microservices/auth/service"
	"github.com/ranefattesingh/microservices/pkg/pgx/pool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("fail to load config", zap.Error(err))
	}

	pool, err := pool.NewConnectionPool(config.Database.EncodedConnectionString())
	if err != nil {
		logger.Fatal("fail to init database", zap.Error(err))
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	authRepository := repository.NewAuthRepository(pool, redisClient)

	// Create token manager with JWT config
	tokenManager := service.NewTokenManager(
		config.JWT.Secret,
		config.JWT.AccessTokenTTL,
		config.JWT.RefreshTokenTTL,
	)

	authService := service.NewAuthService(authRepository, tokenManager, config.JWT.RefreshTokenTTL)
	authHandler := handler.NewAuthHandler(logger, authService)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	server := http.NewHTTPServer(logger, config.Server)

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
