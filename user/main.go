package main

import (
	"net/http"

	"github.com/ranefattesingh/microservices/user/handler"
	"github.com/ranefattesingh/microservices/user/repository"
	"github.com/ranefattesingh/microservices/user/router"
	"github.com/ranefattesingh/microservices/user/service"
)

func main() {
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	router := router.NewRouter(userHandler)
	http.ListenAndServe(":8080", router)
}
