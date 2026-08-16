package main

// @title User Service API
// @version 1.0
// @description Microservice for user management.
// @host localhost:8080
// @BasePath /
import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ranefattesingh/microservices/user/config"
	"github.com/ranefattesingh/microservices/user/handler"
	"github.com/ranefattesingh/microservices/user/repository/db"
	"github.com/ranefattesingh/microservices/user/repository/db/pool"
	"github.com/ranefattesingh/microservices/user/server/http"
	"github.com/ranefattesingh/microservices/user/service"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config, err := config.DefaultLoader().Load()
	if err != nil {
		logger.Fatal("fail to load config", zap.Error(err))
	}

	pool, err := pool.NewConnectionPool(config.Database.EncodedConnectionString())
	if err != nil {
		logger.Fatal("fail to init database", zap.Error(err))
	}

	userRepository := db.NewUserRepository(pool)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(logger, userService)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	server := http.NewHTTPServer(logger, config.Server)

	go func() {
		if err := server.Start(userHandler); err != nil {
			logger.Error("HTTP server failed", zap.Error(err))
			stop()
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(); err != nil {
		logger.Error("HTTP server shutdown failed", zap.Error(err))
	}
}
