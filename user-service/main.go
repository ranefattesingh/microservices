package main

import (
	"github.com/ranefattesingh/ecommerce-platform/user-service/handler"
	"github.com/ranefattesingh/ecommerce-platform/user-service/repository"
	"github.com/ranefattesingh/ecommerce-platform/user-service/repository/postgres/connection"
	"github.com/ranefattesingh/ecommerce-platform/user-service/server"
	"github.com/ranefattesingh/ecommerce-platform/user-service/service"
)

const connStr = "postgres://postgres:postgres@localhost:5432/users?sslmode=disable"

func main() {
	connection, err := connection.NewConnectionPool(connStr)
	if err != nil {
		panic("FAIL")
	}

	accountsRepository := repository.NewAccountsRepository(connection)
	accountsService := service.NewAccountsService(accountsRepository)
	accountsHandler := handler.NewAccountsHandler(accountsService)
	srv := server.NewHTTPServer(accountsHandler)
	srv.Start()
}
