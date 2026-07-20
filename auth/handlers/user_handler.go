package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/ecommerce-platform/auth/handlers/models"
	"github.com/ranefattesingh/ecommerce-platform/auth/router"
	"github.com/ranefattesingh/ecommerce-platform/auth/service"
	"github.com/ranefattesingh/ecommerce-platform/pkg/httperror"
	"github.com/ranefattesingh/ecommerce-platform/pkg/response"
	"go.uber.org/zap"
)

type UsersHandler interface {
	Register(c *gin.Context) response.Responder
	Login(c *gin.Context) response.Responder
}

type userHandler struct {
	validate    *validator.Validate
	userService service.AuthService
	logger      *zap.Logger
}

func NewUserHandler(v *validator.Validate, us service.AuthService, l *zap.Logger) UsersHandler {
	return &userHandler{
		validate:    v,
		userService: us,
		logger:      l,
	}
}

func (uh userHandler) Routes() router.RouterGroup {
	return router.RouterGroup{
		Name: "auth",
		Routes: []router.Route{
			{
				Name:        "Register",
				Path:        "/register",
				Method:      http.MethodPost,
				HandlerFunc: uh.Register,
			},
			{
				Name:        "Login",
				Path:        "/login",
				Method:      http.MethodPost,
				HandlerFunc: uh.Login,
			},
		},
	}
}

func (uh userHandler) Register(c *gin.Context) response.Responder {
	var req models.RegisterRequest
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

	_, err = uh.userService.Register(c, req)
	if err != nil {
		uh.logger.Error("server processing error", zap.Error(err))

		return httperror.InternalServerError()
	}

	return response.NoContent()
}

func (uh userHandler) Login(c *gin.Context) response.Responder {
	var req models.LoginRequest
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

	err = uh.userService.Login(c, req)
	if err != nil {
		uh.logger.Error("invalid user_id param received", zap.Error(err))

		voilations := httperror.Violations{}
		voilations.Add("id", "invalid or empty user id")

		return httperror.BadRequest(voilations)
	}

	return response.NoContent()
}
