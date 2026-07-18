package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ranefattesingh/ecommerce-platform/users/internal/user/handlers/models"
	dbModel "github.com/ranefattesingh/ecommerce-platform/users/internal/user/repository/models"
	"github.com/ranefattesingh/ecommerce-platform/users/internal/user/repository/psql"
	httperror "github.com/ranefattesingh/ecommerce-platform/users/pkg/httperror"
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

	voilations := make(httperror.Violations, 0)

	existingUser, err := s.usersRepo.GetUserHavingEmailOrPhone(ctx, user.Email, user.Phone)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("service:CreateUser {%w}", err)
	}

	if existingUser.Email == user.Email {
		voilations.Add("email", "Email address already taken")
	}

	if existingUser.Phone == user.Phone {
		voilations.Add("phone", "Phone number already taken")
	}

	if voilations.Len() > 0 {
		return 0, httperror.Conflict(voilations)
	}

	id, err := s.usersRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}
