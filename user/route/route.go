package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ranefattesingh/microservices/user/core"
	"github.com/ranefattesingh/microservices/user/handler"
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

	g := gin.New()
	g.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	return &Gin{g}
}

func (g *Gin) Handle(svc core.UserService) {
	ph := handler.NewPingHandler()
	uh := handler.NewUserHandler(svc)

	g.GET("/ping", ph.Ping)

	v1ApiGroup := g.Group("/api/" + apiVersion)

	v1ApiGroup.POST("/login", uh.LoginUser)

	authRoute := v1ApiGroup.Group("/auth", uh.AuthorizeUser())

	// GET
	authRoute.GET("/users", uh.GetAllUsers)
	authRoute.GET("/users/:id", uh.GetUser)

	//POST
	authRoute.POST("/users", uh.AddUser)

	//PUT
	authRoute.PUT("/users/:id", uh.EditUser)

}
