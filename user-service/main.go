package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranefattesingh/ecommerce-platform/user-service/handler"
	"github.com/ranefattesingh/ecommerce-platform/user-service/repository"
	"github.com/ranefattesingh/ecommerce-platform/user-service/service"
)

func main() {
	r := chi.NewRouter()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewHandler(userService)

	r.Post("/users/register", userHandler.RegisterUser)

	http.ListenAndServe(":8000", r)
}
