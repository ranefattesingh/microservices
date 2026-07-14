package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router interface {
	Routes() RouterGroup
}

type RouterGroup struct {
	Name   string
	Routes []Route
}

type Route struct {
	Name        string
	Path        string
	Method      string
	HandlerFunc func(*gin.Context)
}

func NewV1Router(routers ...Router) *gin.Engine {
	mux := gin.New()

	mux.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	mux.Use(gin.Logger())
	mux.Use(gin.Recovery())

	v1Group := mux.Group("/api/v1")

	for _, router := range routers {
		group := v1Group.Group(router.Routes().Name)
		for _, route := range router.Routes().Routes {
			group.Handle(route.Method, route.Path, route.HandlerFunc)
		}
	}

	return mux
}
