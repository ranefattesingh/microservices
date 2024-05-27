package route

import (
	"github.com/gin-gonic/gin"

	"github.com/ranefattesingh/microservices/user/api"
	"github.com/ranefattesingh/microservices/user/core"
	"github.com/ranefattesingh/microservices/user/middleware"
)

const (
	apiVersion = "v1"
)

type Gin struct {
	*gin.Engine
}

func NewRouter(mode string) *Gin {
	gin.SetMode(mode)

	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	return &Gin{router}
}

func (g *Gin) Handle(svc core.UserService) {
	pingHandler := api.NewPingHandler()
	userHandler := api.NewUserHandler(svc)

	g.GET("/ping", pingHandler.Ping)

	v1ApiGroup := g.Group("/api/" + apiVersion)

	v1ApiGroup.POST("/login", userHandler.LoginUser)

	// Auth Routes
	authRoute := v1ApiGroup.Group("/auth", userHandler.AuthorizeUser())
	authRoute.GET("/users", userHandler.RetrieveAllUsers)
	authRoute.GET("/users/:id", userHandler.RetrieveUser)
	authRoute.POST("/users", userHandler.AddUser)
	authRoute.PUT("/users/:id", userHandler.EditUser)
}
