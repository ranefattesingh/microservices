package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ranefattesingh/microservices/auth/handler/dto"
	"github.com/ranefattesingh/microservices/auth/service"
	"go.uber.org/zap"
)

// MockAuthService is a mock implementation of AuthService for testing
type MockAuthService struct {
	loginFunc    func(ctx context.Context, l dto.LoginRequest) (string, string, error)
	registerFunc func(ctx context.Context, r dto.RegisterRequest) error
	refreshFunc  func(ctx context.Context, r dto.RefreshRequest) (string, error)
}

func (m *MockAuthService) Login(ctx context.Context, l dto.LoginRequest) (string, string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, l)
	}
	return "access", "refresh", nil
}

func (m *MockAuthService) Register(ctx context.Context, r dto.RegisterRequest) error {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, r)
	}
	return nil
}

func (m *MockAuthService) Refresh(ctx context.Context, r dto.RefreshRequest) (string, error) {
	if m.refreshFunc != nil {
		return m.refreshFunc(ctx, r)
	}
	return "new-token", nil
}

// TestRegisterHandler tests the register handler
func TestRegisterHandler(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockService := &MockAuthService{
		registerFunc: func(ctx context.Context, r dto.RegisterRequest) error {
			return nil
		},
	}

	handler := NewAuthHandler(logger, mockService)

	body := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

// TestRegisterHandlerBadRequest tests register handler with invalid request
func TestRegisterHandlerBadRequest(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockService := &MockAuthService{}

	handler := NewAuthHandler(logger, mockService)

	req := httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestRegisterHandlerDuplicateEmail tests register handler with duplicate email
func TestRegisterHandlerDuplicateEmail(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockService := &MockAuthService{
		registerFunc: func(ctx context.Context, r dto.RegisterRequest) error {
			return service.ErrEmailAlreadyTaken
		},
	}

	handler := NewAuthHandler(logger, mockService)

	body := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Register(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

// TestLoginHandler tests the login handler
func TestLoginHandler(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockService := &MockAuthService{
		loginFunc: func(ctx context.Context, l dto.LoginRequest) (string, string, error) {
			return "access-token", "refresh-token", nil
		},
	}

	handler := NewAuthHandler(logger, mockService)

	body := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Login(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["token"] != "access-token" {
		t.Fatalf("expected token in response")
	}

	if response["refresh_token"] != "refresh-token" {
		t.Fatalf("expected refresh_token in response")
	}
}

// TestLoginHandlerInvalidCredentials tests login with invalid credentials
func TestLoginHandlerInvalidCredentials(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockService := &MockAuthService{
		loginFunc: func(ctx context.Context, l dto.LoginRequest) (string, string, error) {
			return "", "", service.ErrInvalidCredentials
		},
	}

	handler := NewAuthHandler(logger, mockService)

	body := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestRefreshHandler tests the refresh handler
func TestRefreshHandler(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockService := &MockAuthService{
		refreshFunc: func(ctx context.Context, r dto.RefreshRequest) (string, error) {
			return "new-access-token", nil
		},
	}

	handler := NewAuthHandler(logger, mockService)

	body := dto.RefreshRequest{
		RefreshToken: "refresh-token",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/auth/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Refresh(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["token"] != "new-access-token" {
		t.Fatalf("expected token in response")
	}
}

// TestRefreshHandlerInvalidToken tests refresh with invalid token
func TestRefreshHandlerInvalidToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockService := &MockAuthService{
		refreshFunc: func(ctx context.Context, r dto.RefreshRequest) (string, error) {
			return "", service.ErrInvalidCredentials
		},
	}

	handler := NewAuthHandler(logger, mockService)

	body := dto.RefreshRequest{
		RefreshToken: "invalid-token",
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/v1/auth/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.Refresh(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestHandleError tests the error handler
func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "BadRequest",
			err:            ErrBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "InvalidCredentials",
			err:            service.ErrInvalidCredentials,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "EmailAlreadyTaken",
			err:            service.ErrEmailAlreadyTaken,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "UnknownError",
			err:            errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handleError(w, tt.err)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
