package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ranefattesingh/ecommerce-platform/pkg/httperror"
	"github.com/ranefattesingh/ecommerce-platform/pkg/response"
	"github.com/ranefattesingh/ecommerce-platform/users/handlers/models"
	"github.com/ranefattesingh/ecommerce-platform/users/router"
	"github.com/ranefattesingh/ecommerce-platform/users/service"
	"go.uber.org/zap"
)

type UsersHandler interface {
	CreateUser(c *gin.Context) response.Responder
	GetUser(c *gin.Context) response.Responder
	UpdateUser(c *gin.Context) response.Responder
	DeleteUser(c *gin.Context) response.Responder
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
			{
				Name:        "GetUser",
				Path:        "/:id",
				Method:      http.MethodGet,
				HandlerFunc: uh.GetUser,
			},
			{
				Name:        "UpdateUser",
				Path:        "/:id",
				Method:      http.MethodPut,
				HandlerFunc: uh.UpdateUser,
			},
			{
				Name:        "DeleteUser",
				Path:        "/:id",
				Method:      http.MethodDelete,
				HandlerFunc: uh.DeleteUser,
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

		return handleServiceError(err)
	}

	result := map[string]int64{
		"id": userID,
	}

	location := "/users/" + strconv.FormatInt(userID, 10)

	return response.Created(location).Body(result)
}

func (uh userHandler) GetUser(c *gin.Context) response.Responder {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		uh.logger.Error("invalid user_id param received", zap.Error(err))

		voilations := httperror.Violations{}
		voilations.Add("id", "invalid or empty user id")

		return httperror.BadRequest(voilations)
	}

	user, err := uh.userService.GetUser(c, id)

	if err != nil {
		uh.logger.Error("failed to get user", zap.Error(err))

		return handleServiceError(err)
	}

	return response.Ok(user)
}

func (uh userHandler) UpdateUser(c *gin.Context) response.Responder {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		uh.logger.Error("invalid user_id param received", zap.Error(err))

		voilations := httperror.Violations{}
		voilations.Add("id", "invalid or empty user id")

		return httperror.BadRequest(voilations)
	}

	var req models.UpdateUserRequest
	err = c.BindJSON(&req)
	if err != nil {
		uh.logger.Error("fail to bind the request", zap.Error(err))

		return httperror.BadRequest()
	}

	err = uh.validate.Struct(req)
	if err != nil {
		uh.logger.Error("request validation fail", zap.Error(err))

		return httperror.BadRequest()
	}

	err = uh.userService.UpdateUser(c, id, req)
	if err != nil {
		uh.logger.Error("failed to update user", zap.Error(err))

		return handleServiceError(err)
	}

	return response.NoContent()
}

func (uh userHandler) DeleteUser(c *gin.Context) response.Responder {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		uh.logger.Error("invalid user_id param received", zap.Error(err))

		voilations := httperror.Violations{}
		voilations.Add("id", "invalid or empty user id")

		return httperror.BadRequest(voilations)
	}

	err = uh.userService.DeleteUser(c, id)
	if err != nil {
		uh.logger.Error("failed to delete user", zap.Error(err))

		return handleServiceError(err)
	}

	return response.NoContent()
}

func handleServiceError(err error) *httperror.ProblemDetails {
	if svcErr, ok := errors.AsType[*service.Error](err); ok {
		var violations httperror.Violations
		if len(svcErr.Violations) > 0 {
			violations = make(httperror.Violations, 0)
			for _, v := range svcErr.Violations {
				violations.Add(v.Field, v.Message)
			}
		}

		switch svcErr.Code {
		case service.Validation:
			return httperror.BadRequest(violations)

		case service.NotFound:
			notfound := httperror.NotFound()
			notfound.Detail = svcErr.Message

			return notfound

		case service.Conflict:
			return httperror.Conflict(violations)
		}
	}

	return httperror.InternalServerError()
}
