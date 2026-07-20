package repository

import (
	"context"
	"errors"

	"github.com/ranefattesingh/ecommerce-platform/users/repository/models"
)

var ErrNoRowsUpdated = errors.New("no rows were updated")

type UsersRepository interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUsersHavingEmailOrPhone(ctx context.Context, email, phone string) ([]models.User, error)
	GetUser(ctx context.Context, id int64) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int64) error
}
