package core

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Login(ctx context.Context, login *Login) (string, error)
	Authorize(s string) (map[string]any, error)
	CreateUser(ctx context.Context, user User) (uuid.UUID, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
	EditUser(ctx context.Context, user *User) error
	GetAllUsers(ctx context.Context) (Users, error)
}
