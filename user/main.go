package main

import (
	"net/http"

	"github.com/ranefattesingh/microservices/user/handler"
	"github.com/ranefattesingh/microservices/user/repository/db"
	"github.com/ranefattesingh/microservices/user/repository/db/pool"
	"github.com/ranefattesingh/microservices/user/router"
	"github.com/ranefattesingh/microservices/user/service"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	pool, err := pool.NewConnectionPool("")
	if err != nil {

	}

	userRepository := db.NewUserRepository(pool)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	router := router.NewRouter(userHandler)

	logger.Info("starting server", zap.Int("port", 8080))
	http.ListenAndServe(":8080", router)
}
