package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ranefattesingh/ecommerce-platform/pkg/response"
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
	HandlerFunc HandlerFunc
}

type HandlerFunc func(*gin.Context) response.Responder

func wrap(fn HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := fn(c)
		if resp == nil {
			return
		}

		fmt.Println(resp)

		if err := resp.Respond(c); err != nil {

		}
	}
}

func NewV1Router(routers ...Router) *gin.Engine {
	mux := gin.New()

	mux.GET("/ping", wrap(func(ctx *gin.Context) response.Responder {
		return response.Ok(gin.H{"message": "pong"})
	}))

	mux.Use(gin.Logger())
	mux.Use(gin.Recovery())

	v1Group := mux.Group("/api/v1")

	for _, router := range routers {
		group := v1Group.Group(router.Routes().Name)
		for _, route := range router.Routes().Routes {
			group.Handle(route.Method, route.Path, wrap(route.HandlerFunc))
		}
	}

	return mux
}
