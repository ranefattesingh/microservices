package repository

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/user-service/models"
)

type userRepository struct {
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.RegisterRequest) error
}

func NewUserRepository() *userRepository {
	return &userRepository{}
}

func (r *userRepository) CreateUser(ctx context.Context, createUser models.RegisterRequest) error {
	return nil
}
