package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *userService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) Authorize(s string) (map[string]any, error) {
	return nil, nil
}

func (u *userService) GetAllUsers(ctx context.Context) (Users, error) {
	return nil, nil
}

func (u *userService) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return user, fmt.Errorf("service{GetUserById}: %w", err)
	}

	return user, nil
}

func (u *userService) CreateUser(ctx context.Context, user User) (uuid.UUID, error) {
	userID, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("service{CreateUser}: %w", err)
	}

	return userID, nil
}

func (u *userService) EditUser(ctx context.Context, user *User) error {
	return nil
}

func (u *userService) Login(ctx context.Context, login *Login) (string, error) {
	return "loggedin", nil
}
