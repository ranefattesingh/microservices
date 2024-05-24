package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ranefattesingh/microservices/user/handler"
	"github.com/ranefattesingh/microservices/user/middleware"
)

type Gin struct {
	*gin.Engine
}

func NewRouter(mode string) *Gin {
	gin.SetMode(mode)

	g := gin.New()
	g.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	return &Gin{g}
}

func (g *Gin) Handle() {
	ph := handler.NewPingHandler()

	g.GET("/ping", ph.Ping)
}
