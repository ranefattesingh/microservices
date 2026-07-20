package service

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/auth/handlers/models"
	"github.com/ranefattesingh/ecommerce-platform/auth/repository"
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (int64, error)
	Login(ctx context.Context, req models.LoginRequest) error
}

type authService struct {
	authRepo repository.UsersRepository
}

func NewAuthService(authRepo repository.UsersRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (int64, error) {
	panic("implement me")
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) error {
	panic("implement me")
}
