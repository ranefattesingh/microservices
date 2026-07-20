package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	handlerspkg "github.com/ranefattesingh/ecommerce-platform/users/handlers"
	handlersmodels "github.com/ranefattesingh/ecommerce-platform/users/handlers/models"
	"github.com/ranefattesingh/ecommerce-platform/users/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestContext(t *testing.T, method, target string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(method, target, strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx, w
}

func TestCreateUser_Success(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	mockService.On("CreateUser", mock.Anything, mock.MatchedBy(func(req handlersmodels.CreateUserRequest) bool {
		return req.Email == "ada@example.com"
	})).Return(int64(55), nil).Once()

	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, w := newTestContext(t, http.MethodPost, "/users", `{
		"firstName":"Ada",
		"lastName":"Lovelace",
		"email":"ada@example.com",
		"phone":"1234567890",
		"password":"secret123",
		"confirmPassword":"secret123",
		"accessType":"admin"
	}`)

	responder := handler.CreateUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "/users/55", w.Header().Get("Location"))
	mockService.AssertExpectations(t)
}

func TestCreateUser_ReturnsBadRequestForValidationFailure(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, w := newTestContext(t, http.MethodPost, "/users", `{"firstName":"Ad"}`)

	responder := handler.CreateUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
}

func TestGetUser_Success(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	mockService.On("GetUser", mock.Anything, int64(7)).Return(handlersmodels.User{ID: 7, Email: "grace@example.com"}, nil).Once()

	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, w := newTestContext(t, http.MethodGet, "/users/7", "")
	ctx.Params = gin.Params{{Key: "id", Value: "7"}}

	responder := handler.GetUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetUser_ReturnsBadRequestForInvalidId(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, w := newTestContext(t, http.MethodGet, "/users/not-a-number", "")
	ctx.Params = gin.Params{{Key: "id", Value: "not-a-number"}}

	responder := handler.GetUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "GetUser", mock.Anything, mock.Anything)
}

func TestUpdateUser_Success(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	mockService.On("UpdateUser", mock.Anything, int64(11), mock.Anything).Return(nil).Once()

	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, _ := newTestContext(t, http.MethodPut, "/users/11", `{
		"firstName":"Linus",
		"lastName":"Torvalds",
		"email":"linus@example.com",
		"phone":"1234567890"
	}`)
	ctx.Params = gin.Params{{Key: "id", Value: "11"}}
	responder := handler.UpdateUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())
	mockService.AssertExpectations(t)
}

func TestDeleteUser_MapsServiceErrorsToNotFound(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	mockService.On("DeleteUser", mock.Anything, int64(99)).Return(&service.Error{Code: service.NotFound, Message: "User profile not found"}).Once()

	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, w := newTestContext(t, http.MethodDelete, "/users/99", "")
	ctx.Params = gin.Params{{Key: "id", Value: "99"}}

	responder := handler.DeleteUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandleServiceError_MapsConflictToConflictStatus(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	mockService.On("GetUser", mock.Anything, int64(3)).Return(handlersmodels.User{}, &service.Error{Code: service.Conflict, Message: "conflict"}).Once()

	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, w := newTestContext(t, http.MethodGet, "/users/3", "")
	ctx.Params = gin.Params{{Key: "id", Value: "3"}}

	responder := handler.GetUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandleServiceError_MapsValidationErrorToBadRequest(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	mockService.On("GetUser", mock.Anything, int64(4)).Return(handlersmodels.User{}, &service.Error{Code: service.Validation, Message: "validation"}).Once()

	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, w := newTestContext(t, http.MethodGet, "/users/4", "")
	ctx.Params = gin.Params{{Key: "id", Value: "4"}}

	responder := handler.GetUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandleServiceError_DefaultsToInternalServerError(t *testing.T) {
	t.Parallel()

	mockService := service.NewMockUsersService()
	mockService.On("GetUser", mock.Anything, int64(5)).Return(handlersmodels.User{}, errors.New("boom")).Once()

	handler := handlerspkg.NewUserHandler(validator.New(), mockService, zap.NewNop())
	ctx, w := newTestContext(t, http.MethodGet, "/users/5", "")
	ctx.Params = gin.Params{{Key: "id", Value: "5"}}

	responder := handler.GetUser(ctx)
	require.NoError(t, responder.Respond(ctx))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
