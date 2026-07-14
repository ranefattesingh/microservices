package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/handlers/models"
	dbModel "github.com/ranefattesingh/ecommerce-platform/internal/user/repository/models"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/repository/psql"
	httperror "github.com/ranefattesingh/ecommerce-platform/pkg/http_error"
)

var (
	ErrEmailAlreadyTaken = &httperror.ErrorInfo{
		Message:        "email is already taken",
		HTTPStatusCode: http.StatusConflict,
	}

	ErrPhoneAlreadyTaken = &httperror.ErrorInfo{
		Message:        "phone is already taken",
		HTTPStatusCode: http.StatusConflict,
	}
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

	existingUser, err := s.usersRepo.GetUserHavingEmailOrPhone(ctx, user.Email, user.Phone)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("service:CreateUser {%w}", err)
	}

	if existingUser.Email == user.Email {
		return 0, fmt.Errorf("service:CreateUser {%w}", ErrEmailAlreadyTaken)
	}

	if existingUser.Phone == user.Phone {
		return 0, fmt.Errorf("service:CreateUser {%w}", ErrPhoneAlreadyTaken)
	}

	id, err := s.usersRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("service:CreateUser {%w}", err)
	}

	return id, nil
}
