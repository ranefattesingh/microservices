package service

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/user-service/models"
	"github.com/ranefattesingh/ecommerce-platform/user-service/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, createUser models.RegisterRequest) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) CreateUser(ctx context.Context, createUser models.RegisterRequest) error {
	return s.userRepository.CreateUser(ctx, createUser)
}
