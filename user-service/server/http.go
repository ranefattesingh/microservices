package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ranefattesingh/ecommerce-platform/user-service/handler"
)

const (
	addr    = ":8080"
	baseURL = "/api/v1"
)

type HTTPServer struct {
	accountHandler handler.AccountsHandler
}

func NewHTTPServer(ah handler.AccountsHandler) *HTTPServer {
	return &HTTPServer{
		accountHandler: ah,
	}
}

func (s *HTTPServer) Start() error {
	g := gin.New()

	v1Group := g.Group(baseURL)
	v1Group.POST("/accounts", s.accountHandler.CreateAccount)

	err := g.Run(addr)
	if err != nil {
		return fmt.Errorf("server: {%w}", err)
	}

	return nil
}
