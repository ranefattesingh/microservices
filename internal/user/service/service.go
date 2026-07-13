package service

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/internal/user/handlers/models"
)

type UsersService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (int32, error)
}

type usersService struct{}

func (us *usersService) CreateUser(ctx context.Context, req models.CreateUserRequest) (int32, error) {
	return 1, nil
}
