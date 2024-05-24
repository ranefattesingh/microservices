package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler interface {
	Ping(c *gin.Context)
}

type pingHandle struct{}

func NewPingHandler() *pingHandle {
	return new(pingHandle)
}

func (ph *pingHandle) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong!")
}
