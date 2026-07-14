package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/ecommerce-platform/internal/router"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/handlers/models"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/service"
	"github.com/ranefattesingh/ecommerce-platform/pkg/response"
	"go.uber.org/zap"
)

type UsersHandler interface {
	CreateUser(c *gin.Context)
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

func (uh userHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		uh.logger.Error("fail to bind the request", zap.Error(err))

		response.BadRequest(c, err)

		return
	}

	err = uh.validate.Struct(req)
	if err != nil {
		uh.logger.Error("request validation fail", zap.Error(err))
		response.BadRequest(c, err)

		return
	}

	resp, err := uh.userService.CreateUser(c, req)
	if err != nil {
		uh.logger.Error("server processing error", zap.Error(err))
		response.ErrorResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, map[string]int64{"id": resp})
}
