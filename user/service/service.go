package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ranefattesingh/microservices/user/pkg"
	"github.com/ranefattesingh/microservices/user/repository/db"
	dbModel "github.com/ranefattesingh/microservices/user/repository/db/models"
	"github.com/ranefattesingh/microservices/user/service/models"
)

var (
	ErrEmailAlreadyTaken = errors.New("user with given email already present")
	ErrUserDoesNotExist  = errors.New("user does not exist")
)
var _ UserService = (*userService)(nil)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, user *models.User) error
}

type userService struct {
	repo db.UserRepository
}

func NewUserService(repo db.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	_, err := s.repo.GetByEmail(ctx, user.Email)
	if err == nil {
		return -1, ErrEmailAlreadyTaken
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return -1, fmt.Errorf("service: %w", err)
	}

	hashed, err := pkg.EncryptPassword(user.Password)
	if err != nil {
		return -1, fmt.Errorf("service: %w", err)
	}

	model := &dbModel.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashed,
	}

	id, err := s.repo.Create(ctx, model)
	if err != nil {
		return 0, fmt.Errorf("service: %w", err)
	}

	return id, nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserDoesNotExist
		}

		return nil, fmt.Errorf("service: %w", err)
	}

	model := &models.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return model, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, updateUser *models.User) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserDoesNotExist
		}

		return fmt.Errorf("service: %w", err)
	}

	// attempted to updated user email with different user
	userWithEmail, err := s.repo.GetByEmail(ctx, updateUser.Email)
	if err == nil &&
		updateUser.Email == userWithEmail.Email &&
		id != userWithEmail.ID {
		fmt.Println(id)
		fmt.Println(userWithEmail)
		return ErrEmailAlreadyTaken
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("service: %w", err)
	}

	user.FirstName = updateUser.FirstName
	user.LastName = updateUser.LastName
	user.Email = updateUser.Email

	if err := s.repo.Update(ctx, id, user); err != nil {
		return fmt.Errorf("service: %w", err)
	}

	return nil
}
