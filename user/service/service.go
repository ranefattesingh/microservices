package service

import (
	"context"

	"github.com/ranefattesingh/microservices/user/repository"
	"github.com/ranefattesingh/microservices/user/service/models"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
