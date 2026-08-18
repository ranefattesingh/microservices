package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ranefattesingh/microservices/auth/handler/dto"
	"github.com/ranefattesingh/microservices/auth/repository"
	"github.com/ranefattesingh/microservices/auth/repository/db/models"
)

// MockAuthRepository is a mock implementation of AuthRepository for testing
type MockAuthRepository struct {
	accounts         map[string]*models.Account
	refreshTokens    map[string]string
	createAccountErr error
	getByEmailErr    error
}

// NewMockAuthRepository creates a new mock repository
func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		accounts:      make(map[string]*models.Account),
		refreshTokens: make(map[string]string),
	}
}

// CreateAccount creates an account
func (m *MockAuthRepository) CreateAccount(ctx context.Context, account *models.Account) (int64, error) {
	if m.createAccountErr != nil {
		return 0, m.createAccountErr
	}

	if _, exists := m.accounts[account.Email]; exists {
		return 0, repository.ErrEmailAlreadyExists
	}

	account.ID = int64(len(m.accounts) + 1)
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	m.accounts[account.Email] = account
	return account.ID, nil
}

// GetAccountByEmail gets an account by email
func (m *MockAuthRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}

	account, exists := m.accounts[email]
	if !exists {
		return nil, repository.ErrUserNotFound
	}

	return account, nil
}

// SaveRefreshToken saves a refresh token
func (m *MockAuthRepository) SaveRefreshToken(ctx context.Context, token, email string, ttl time.Duration) error {
	m.refreshTokens[token] = email
	return nil
}

// GetEmailByRefreshToken gets email by refresh token
func (m *MockAuthRepository) GetEmailByRefreshToken(ctx context.Context, token string) (string, error) {
	email, exists := m.refreshTokens[token]
	if !exists {
		return "", errors.New("token not found")
	}
	return email, nil
}

// TestRegisterSuccess tests successful registration
func TestRegisterSuccess(t *testing.T) {
	mockRepo := NewMockAuthRepository()
	tokenManager := NewTokenManager("secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockRepo, tokenManager, 7*24*time.Hour)

	req := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	err := svc.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify account was created
	account, err := mockRepo.GetAccountByEmail(context.Background(), req.Email)
	if err != nil {
		t.Fatalf("expected account to exist, got error: %v", err)
	}

	if account.Email != req.Email {
		t.Fatalf("expected email %s, got %s", req.Email, account.Email)
	}

	// Verify password is hashed (not plain text)
	if account.Password == req.Password {
		t.Fatalf("expected password to be hashed")
	}
}

// TestRegisterDuplicateEmail tests duplicate email registration
func TestRegisterDuplicateEmail(t *testing.T) {
	mockRepo := NewMockAuthRepository()
	tokenManager := NewTokenManager("secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockRepo, tokenManager, 7*24*time.Hour)

	req := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// First registration should succeed
	err := svc.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	// Second registration with same email should fail
	err = svc.Register(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error for duplicate email, got nil")
	}

	if !errors.Is(err, ErrEmailAlreadyTaken) {
		t.Fatalf("expected ErrEmailAlreadyTaken, got %v", err)
	}
}

// TestRegisterShortPassword tests registration with short password
func TestRegisterShortPassword(t *testing.T) {
	mockRepo := NewMockAuthRepository()
	tokenManager := NewTokenManager("secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockRepo, tokenManager, 7*24*time.Hour)

	req := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "short",
	}

	err := svc.Register(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error for short password, got nil")
	}
}

// TestLoginSuccess tests successful login
func TestLoginSuccess(t *testing.T) {
	mockRepo := NewMockAuthRepository()
	tokenManager := NewTokenManager("secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockRepo, tokenManager, 7*24*time.Hour)

	// Register first
	regReq := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	err := svc.Register(context.Background(), regReq)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	// Login
	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	accessToken, refreshToken, err := svc.Login(context.Background(), loginReq)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if accessToken == "" {
		t.Fatalf("expected access token, got empty string")
	}

	if refreshToken == "" {
		t.Fatalf("expected refresh token, got empty string")
	}

	// Verify tokens are valid JWTs
	claims, err := tokenManager.ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("expected valid access token, got error: %v", err)
	}

	if claims.Email != "test@example.com" {
		t.Fatalf("expected email in token, got %s", claims.Email)
	}
}

// TestLoginInvalidCredentials tests login with wrong credentials
func TestLoginInvalidCredentials(t *testing.T) {
	mockRepo := NewMockAuthRepository()
	tokenManager := NewTokenManager("secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockRepo, tokenManager, 7*24*time.Hour)

	// Register first
	regReq := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	svc.Register(context.Background(), regReq)

	// Try login with wrong password
	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, _, err := svc.Login(context.Background(), loginReq)
	if err == nil {
		t.Fatalf("expected error for invalid credentials, got nil")
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

// TestLoginNonexistentUser tests login for non-existent user
func TestLoginNonexistentUser(t *testing.T) {
	mockRepo := NewMockAuthRepository()
	tokenManager := NewTokenManager("secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockRepo, tokenManager, 7*24*time.Hour)

	loginReq := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, _, err := svc.Login(context.Background(), loginReq)
	if err == nil {
		t.Fatalf("expected error for nonexistent user, got nil")
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

// TestRefreshSuccess tests successful token refresh
func TestRefreshSuccess(t *testing.T) {
	mockRepo := NewMockAuthRepository()
	tokenManager := NewTokenManager("secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockRepo, tokenManager, 7*24*time.Hour)

	// Register and login
	regReq := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	svc.Register(context.Background(), regReq)

	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	_, refreshToken, err := svc.Login(context.Background(), loginReq)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// Refresh token
	refreshReq := dto.RefreshRequest{
		RefreshToken: refreshToken,
	}

	newAccessToken, err := svc.Refresh(context.Background(), refreshReq)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newAccessToken == "" {
		t.Fatalf("expected new access token, got empty string")
	}

	// Verify new token is valid
	claims, err := tokenManager.ValidateAccessToken(newAccessToken)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}

	if claims.Email != "test@example.com" {
		t.Fatalf("expected email in token, got %s", claims.Email)
	}
}

// TestRefreshInvalidToken tests refresh with invalid token
func TestRefreshInvalidToken(t *testing.T) {
	mockRepo := NewMockAuthRepository()
	tokenManager := NewTokenManager("secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewAuthService(mockRepo, tokenManager, 7*24*time.Hour)

	refreshReq := dto.RefreshRequest{
		RefreshToken: "invalid.token.here",
	}

	_, err := svc.Refresh(context.Background(), refreshReq)
	if err == nil {
		t.Fatalf("expected error for invalid token, got nil")
	}

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
