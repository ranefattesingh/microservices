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

var (
	ErrUserNotFound  = httperror.NotFound().SetDetail("user could not be found")
	ErrNoUserUpdated = errors.New("user details could not be updated")
)

type UsersService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (int64, error)
	GetUser(ctx context.Context, id int64) (models.User, error)
	UpdateUser(ctx context.Context, id int64, req models.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id int64) error
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

	existingUsers, err := s.usersRepo.GetUsersHavingEmailOrPhone(ctx, user.Email, user.Phone)
	if err != nil {
		return 0, fmt.Errorf("service:CreateUser {%w}", err)
	}

	if len(existingUsers) > 0 {
		violations := make(httperror.Violations, 0)
		for _, existing := range existingUsers {
			if existing.Email == req.Email {
				violations.Add("email", "Email address already taken")
			}

			if existing.Phone == req.Phone {
				violations.Add("phone", "Phone number already taken")
			}
		}

		if violations.Len() > 0 {
			return 0, fmt.Errorf("service:CreateUser {%w}", httperror.Conflict(violations))
		}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("usersService.GetUser: %w", ErrUserNotFound)
		}

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

func (s *usersService) UpdateUser(ctx context.Context, id int64, req models.UpdateUserRequest) error {
	userDb, err := s.usersRepo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("usersService.UpdateUser: %w", ErrUserNotFound)
		}

		return fmt.Errorf("usersService.UpdateUser: %w", err)
	}

	existingUsers, err := s.usersRepo.GetUsersHavingEmailOrPhone(ctx, req.Email, req.Phone)
	if err != nil {
		return fmt.Errorf("service:UpdateUser {%w}", err)
	}

	if len(existingUsers) > 0 {
		violations := make(httperror.Violations, 0)
		for _, existing := range existingUsers {
			if existing.ID == userDb.ID {
				continue // Skip checking the user modifying their own row
			}

			if existing.ID != userDb.ID {
				if existing.Email == req.Email {
					violations.Add("email", "Email address already taken")
				}
				if existing.Phone == req.Phone {
					violations.Add("phone", "Phone number already taken")
				}
			}
		}

		if violations.Len() > 0 {
			return fmt.Errorf("service:UpdateUser {%w}", httperror.Conflict(violations))
		}
	}

	user := dbModel.User{
		ID:         userDb.ID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Phone:      req.Phone,
		AccessType: userDb.AccessType,
	}

	err = s.usersRepo.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, psql.ErrNoRowsUpdated) {
			return ErrNoUserUpdated
		}

		return fmt.Errorf("service:UpdateUser {%w}", err)
	}

	return nil
}

func (s *usersService) DeleteUser(ctx context.Context, id int64) error {
	err := s.usersRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, psql.ErrNoRowsUpdated) {
			return fmt.Errorf("service:DeleteUser {%w}", ErrUserNotFound)
		}

		return fmt.Errorf("service:DeleteUser {%w}", err)
	}

	return nil
}
