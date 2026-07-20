package repository

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/auth/repository/models"
)

type UsersRepository interface {
	Register(ctx context.Context, email, passwordHash string) (int64, error)
	Login(ctx context.Context, email, passwordHash string) ([]models.Auth, error)
}
