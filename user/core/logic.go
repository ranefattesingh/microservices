package core

import (
	"context"

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

func (u *userService) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	return nil, nil
}

func (u *userService) CreateUser(ctx context.Context, user User) (uuid.UUID, error) {
	userID, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (u *userService) EditUser(ctx context.Context, user *User) error {
	return nil
}

func (u *userService) Login(ctx context.Context, login *Login) (string, error) {
	return "loggedin", nil
}
