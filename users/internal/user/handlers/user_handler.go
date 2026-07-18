package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/ecommerce-platform/pkg/httperror"
	"github.com/ranefattesingh/ecommerce-platform/pkg/response"
	"github.com/ranefattesingh/ecommerce-platform/users/internal/router"
	"github.com/ranefattesingh/ecommerce-platform/users/internal/user/handlers/models"
	"github.com/ranefattesingh/ecommerce-platform/users/internal/user/service"
	"go.uber.org/zap"
)

type UsersHandler interface {
	CreateUser(c *gin.Context) response.Responder
}

type userHandler struct {
	validate    *validator.Validate
	userService service.UsersService
	logger      *zap.Logger
}

func NewUserHandler(v *validator.Validate, us service.UsersService, l *zap.Logger) *userHandler {
	return &userHandler{
		validate:    v,
		userService: us,
		logger:      l,
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

func (uh userHandler) CreateUser(c *gin.Context) response.Responder {
	var req models.CreateUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		uh.logger.Error("fail to bind the request", zap.Error(err))

		return httperror.BadRequest()
	}

	err = uh.validate.Struct(req)
	if err != nil {
		uh.logger.Error("request validation fail", zap.Error(err))

		return httperror.BadRequest()
	}

	userID, err := uh.userService.CreateUser(c, req)
	if err != nil {
		uh.logger.Error("server processing error", zap.Error(err))

		pb, ok := errors.AsType[*httperror.ProblemDetails](err)
		if ok {
			return pb
		}

		return httperror.InternalServerError(err)
	}

	result := map[string]int64{
		"id": userID,
	}

	location := "/users/" + strconv.FormatInt(userID, 10)

	return response.Created(location).Body(result)
}
