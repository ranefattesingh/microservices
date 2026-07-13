package service

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/internal/user/handlers/models"
	dbModel "github.com/ranefattesingh/ecommerce-platform/internal/user/repository/models"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/repository/psql"
)

type UsersService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (int64, error)
}

type usersService struct {
	usersRepo psql.UsersRepository
}

func NewUsersService(usersRepo psql.UsersRepository) UsersService {
	return &usersService{
		usersRepo: usersRepo,
	}
}

func (s *usersService) CreateUser(ctx context.Context, req models.CreateUserRequest) (int64, error) {
	user := dbModel.User{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Phone:      req.Phone,
		AccessType: dbModel.AccessType(req.AccessType),
	}

	id, err := s.usersRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
