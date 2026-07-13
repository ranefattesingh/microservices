package psql

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/internal/user/repository/models"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, model models.User) (int32, error)
}

type usersRepository struct{}

func (ur *usersRepository) CreateUser(ctx context.Context, model models.User) {

}
