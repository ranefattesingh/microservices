package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ranefattesingh/microservices/user/core"

	"github.com/ranefattesingh/pkg/json"
	"github.com/ranefattesingh/pkg/log"
)

var (
	ErrBadCreateRequest = &json.Error{
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "invalid create user request",
		Code:           http.StatusBadRequest,
	}

	ErrInvalidUserID = &json.Error{
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "invalid user id",
		Code:           http.StatusBadRequest,
	}
)

type UserInterface interface {
	RetrieveUser(ctx *gin.Context)
	RetrieveAllUsers(ctx *gin.Context)
	AddUser(ctx *gin.Context)
	EditUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
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

func (uh *userHandle) RetrieveUser(ctx *gin.Context) {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Logger().Error("user_id parsing error", zap.Error(err))
		logOnError(json.EncodeErrorJSON(ctx.Writer, ErrInvalidUserID))
	}

	user, err := uh.service.GetUserByID(ctx, userID)
	if err != nil {
		log.Logger().Error("retrieve user returned with error", zap.Error(err))
		logOnError(json.EncodeErrorJSON(ctx.Writer, err))
	}

	logOnError(json.EncodeResponseJSON(ctx.Writer, http.StatusOK, user))
}

func (uh *userHandle) RetrieveAllUsers(ctx *gin.Context) {
	users, err := uh.service.GetAllUsers(ctx)
	if err != nil {
		log.Logger().Error("retrieve all users returned with error", zap.Error(err))
		logOnError(json.EncodeErrorJSON(ctx.Writer, err))
	}

	logOnError(json.EncodeResponseJSON(ctx.Writer, http.StatusOK, users))
}

func (uh *userHandle) AddUser(ctx *gin.Context) {
	req, err := json.DecodeAndValidateJSON[core.User](ctx.Request)
	if err != nil {
		var jsonErr *json.Error
		if !errors.As(err, &jsonErr) {
			jsonErr = ErrBadCreateRequest
		}

		log.Logger().Error("request decoding error", zap.Error(err))
		logOnError(json.EncodeErrorJSON(ctx.Writer, jsonErr))

		return
	}

	id, err := uh.service.CreateUser(ctx, req)
	if err != nil {
		log.Logger().Error("add user returned with error", zap.Error(err))
		logOnError(json.EncodeErrorJSON(ctx.Writer, err))

		return
	}

	logOnError(json.EncodeResponseJSON(ctx.Writer, http.StatusCreated, map[string]uuid.UUID{"user_id": id}))
}

func (uh *userHandle) EditUser(ctx *gin.Context) {
	err := uh.service.EditUser(ctx, &core.User{})
	if err != nil {
		log.Logger().Error("edit user returned with error", zap.Error(err))
		logOnError(json.EncodeErrorJSON(ctx.Writer, err))
	}

	logOnError(json.EncodeResponseJSON(ctx.Writer, http.StatusNoContent, nil))
}

func (uh *userHandle) LoginUser(ctx *gin.Context) {
	data, err := uh.service.Login(ctx, &core.Login{})
	if err != nil {
		log.Logger().Error("login user returned with error", zap.Error(err))
		logOnError(json.EncodeErrorJSON(ctx.Writer, err))
	}

	logOnError(json.EncodeResponseJSON(ctx.Writer, http.StatusNoContent, data))
}

func (uh *userHandle) AuthorizeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Logger().Info("invoked AuthorizeUser")

		ctx.Next()
	}
}

func logOnError(err error) {
	if err != nil {
		log.Logger().Error("error while returning encoding json", zap.Error(err))
	}
}
