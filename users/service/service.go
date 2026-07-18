package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	httperror "github.com/ranefattesingh/ecommerce-platform/pkg/httperror"
	"github.com/ranefattesingh/ecommerce-platform/users/handlers/models"
	dbModel "github.com/ranefattesingh/ecommerce-platform/users/repository/models"
	"github.com/ranefattesingh/ecommerce-platform/users/repository/psql"
)

type UsersService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (int64, error)
	GetUser(ctx context.Context, id int64) (models.User, error)
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

func (s *usersService) GetUser(ctx context.Context, id int64) (models.User, error) {
	userDb, err := s.usersRepo.GetUser(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("usersService.GetUser: %w", err)
	}

	user := models.User{
		ID:         userDb.ID,
		FirstName:  userDb.FirstName,
		LastName:   userDb.LastName,
		Email:      userDb.Email,
		Phone:      userDb.Phone,
		AccessType: models.AccessType(userDb.AccessType),
		CreatedAt:  userDb.CreatedAt,
		UpdatedAt:  userDb.UpdatedAt,
	}

	return user, nil
}
