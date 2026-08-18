package service

import (
	"context"
	"errors"
	"time"

	"github.com/ranefattesingh/microservices/auth/handler/dto"
	"github.com/ranefattesingh/microservices/auth/repository"
	"github.com/ranefattesingh/microservices/auth/service/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyTaken  = errors.New("email already taken")
)

type AuthService interface {
	Login(ctx context.Context, l dto.LoginRequest) (string, string, error)
	Register(ctx context.Context, r dto.RegisterRequest) error
	Refresh(ctx context.Context, r dto.RefreshRequest) (string, error)
}

// AuthService provides auth related operations backed by repository.
type authService struct {
	repo            repository.AuthRepository
	tokenManager    *TokenManager
	refreshTokenTTL time.Duration
}

func NewAuthService(repo repository.AuthRepository, tokenManager *TokenManager, refreshTokenTTL time.Duration) AuthService {
	return &authService{
		repo:            repo,
		tokenManager:    tokenManager,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// Register creates a new account in repository.
func (s *authService) Register(ctx context.Context, r dto.RegisterRequest) error {
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// hash password
	hash, err := EncryptPassword(r.Password)
	if err != nil {
		return err
	}

	account := &models.Account{
		Email:    r.Email,
		Password: hash,
	}

	_, err = s.repo.CreateAccount(ctx, account.ToRepositoryModel())
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return ErrEmailAlreadyTaken
		}
		return err
	}

	return nil
}

// Login validates credentials and returns access + refresh tokens.
func (s *authService) Login(ctx context.Context, l dto.LoginRequest) (string, string, error) {
	account, err := s.repo.GetAccountByEmail(ctx, l.Email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if err := ComparePassword(account.Password, l.Password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate JWT tokens
	accessToken, err := s.tokenManager.GenerateAccessToken(l.Email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(l.Email)
	if err != nil {
		return "", "", err
	}

	// store refresh token in repo for revocation support (optional)
	if err := s.repo.SaveRefreshToken(ctx, refreshToken, l.Email, s.refreshTokenTTL); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Refresh validates refresh token via repo then issues new access token.
func (s *authService) Refresh(ctx context.Context, r dto.RefreshRequest) (string, error) {
	// Validate refresh token JWT
	claims, err := s.tokenManager.ValidateRefreshToken(r.RefreshToken)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Also check if refresh token exists in repo (for revocation support)
	_, err = s.repo.GetEmailByRefreshToken(ctx, r.RefreshToken)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate new access token
	newAccessToken, err := s.tokenManager.GenerateAccessToken(claims.Email)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
