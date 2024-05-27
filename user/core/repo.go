package core

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (uuid.UUID, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
}
