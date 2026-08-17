package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	repo "github.com/ranefattesingh/microservices/auth/repository"
	"github.com/ranefattesingh/microservices/auth/repository/cache"
	"github.com/ranefattesingh/microservices/auth/repository/db"
	"github.com/ranefattesingh/microservices/auth/handler/dto"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	defaultRefreshTokenTTL  = 7 * 24 * time.Hour
)

// AuthService provides auth related operations backed by repository.
type AuthService struct {
	repo *repo.UserRepository
}

func NewAuthService() *AuthService {
	ctx := context.Background()
	// init db pool
	dbPool, _ := db.NewPool(ctx)
	// init cache client
	cacheClient := cache.NewClient()

	return &AuthService{
		repo: repo.NewUserRepository(dbPool, cacheClient),
	}
}

// Register creates a new user in repository.
func (s *AuthService) Register(ctx context.Context, r dto.RegisterRequest) error {
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.repo.CreateUser(ctx, r.Email, string(hash))
	if err != nil {
		if errors.Is(err, repo.ErrEmailAlreadyExists) {
			return ErrEmailAlreadyTaken
		}
		return err
	}

	return nil
}

// Login validates credentials and returns access + refresh tokens.
func (s *AuthService) Login(ctx context.Context, l dto.LoginRequest) (string, string, error) {
	_, storedHash, err := s.repo.GetUserByEmail(ctx, l.Email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(l.Password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	access := s.randomToken()
	refresh := s.randomToken()

	// store refresh token in repo
	if err := s.repo.StoreRefreshToken(ctx, refresh, l.Email, defaultRefreshTokenTTL); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// Refresh validates refresh token via repo then issues new access token.
func (s *AuthService) Refresh(ctx context.Context, r dto.RefreshRequest) (string, error) {
	email, err := s.repo.GetEmailByRefreshToken(ctx, r.RefreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// touch or load user from cache (ensures cache is refreshed)
	_, _, _ = s.repo.GetUserByEmail(ctx, email)

	return s.randomToken(), nil
}

func (s *AuthService) randomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}
