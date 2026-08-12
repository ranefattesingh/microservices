package main

import (
	"net/http"

	"github.com/ranefattesingh/microservices/user/handler"
	"github.com/ranefattesingh/microservices/user/router"
)

func main() {
	userHandler := handler.NewUserHandler()
	router := router.NewRouter(userHandler)
	http.ListenAndServe(":8080", router)
}
