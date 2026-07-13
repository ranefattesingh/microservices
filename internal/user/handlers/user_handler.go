package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/handlers/models"
)

type UsersHandler interface {
	CreateUser(c *gin.Context)
}

type userHandler struct {
	validate *validator.Validate
}

func NewUserHandler(v *validator.Validate) *userHandler {
	return &userHandler{validate: v}
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
