package psql

import (
	"context"

	"github.com/ranefattesingh/ecommerce-platform/auth/platform/database/psql"
	"github.com/ranefattesingh/ecommerce-platform/auth/repository/models"
)

type authRepository struct {
	db *psql.Database
}

func NewAuthRepository(db *psql.Database) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(ctx context.Context, email, passwordHash string) (int64, error) {
	panic("implement me")
}

func (r *authRepository) Login(ctx context.Context, email, passwordHash string) ([]models.Auth, error) {
	panic("implement me")
}
