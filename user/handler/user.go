package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranefattesingh/microservices/user/core"
	"github.com/ranefattesingh/microservices/user/models"
	"github.com/ranefattesingh/pkg/json"
	"github.com/ranefattesingh/pkg/log"
)

type UserInterface interface {
	GetUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	AddUser(c *gin.Context)
	EditUser(c *gin.Context)
	LoginUser(c *gin.Context)
	AuthorizeUser() gin.HandlerFunc
}

type userHandle struct {
	service core.UserService
}

func NewUserHandler(svc core.UserService) *userHandle {
	return &userHandle{
		service: svc,
	}
}

func (uh *userHandle) GetUser(c *gin.Context) {
	user, err := uh.service.GetUserById(c, 0)
	if err != nil {
		json.EncodeErrorJSON(c.Writer, err)
	}

	json.EncodeResponseJSON(c.Writer, http.StatusOK, user)
}

func (uh *userHandle) GetAllUsers(c *gin.Context) {
	users, err := uh.service.GetAllUsers(c)
	if err != nil {
		json.EncodeErrorJSON(c.Writer, err)
	}

	json.EncodeResponseJSON(c.Writer, http.StatusOK, users)
}

func (uh *userHandle) AddUser(c *gin.Context) {
	id, err := uh.service.AddUser(c, &models.User{})
	if err != nil {
		json.EncodeErrorJSON(c.Writer, err)
	}

	json.EncodeResponseJSON(c.Writer, http.StatusCreated, id)
}

func (uh *userHandle) EditUser(c *gin.Context) {
	err := uh.service.EditUser(c, &models.User{})
	if err != nil {
		json.EncodeErrorJSON(c.Writer, err)
	}

	json.EncodeResponseJSON(c.Writer, http.StatusNoContent, nil)
}

func (uh *userHandle) LoginUser(c *gin.Context) {
	data, err := uh.service.Login(c, &models.Login{})
	if err != nil {
		json.EncodeErrorJSON(c.Writer, err)
	}

	json.EncodeResponseJSON(c.Writer, http.StatusNoContent, data)
}

func (uh *userHandle) AuthorizeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Logger().Info("invoked AuthorizeUser")

		ctx.Next()
	}
}
