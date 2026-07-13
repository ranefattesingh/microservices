package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/ecommerce-platform/internal/router"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/handlers/models"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/service"
)

type UsersHandler interface {
	CreateUser(c *gin.Context)
}

type userHandler struct {
	validate    *validator.Validate
	userService service.UsersService
}

func NewUserHandler(v *validator.Validate, us service.UsersService) *userHandler {
	return &userHandler{
		validate:    v,
		userService: us,
	}
}

func (uh userHandler) Routes() router.RouterGroup {
	return router.RouterGroup{
		Name: "users",
		Routes: []router.Route{
			{
				Name:        "CreateUser",
				Path:        "",
				Method:      http.MethodPost,
				HandlerFunc: uh.CreateUser,
			},
		},
	}
}

func (uh userHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "fail to bind the request")

		return
	}

	err = uh.validate.Struct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "fail to bind the request")

		return
	}

	// Call Next Layer

	c.JSON(http.StatusOK, "SUCCESS")
}
