package psql

import (
	"context"
	"fmt"

	"github.com/ranefattesingh/ecommerce-platform/internal/platform/database/psql"
	"github.com/ranefattesingh/ecommerce-platform/internal/user/repository/models"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
}

type usersRepository struct {
	db *psql.Database
}

func NewUsersRepository(db *psql.Database) *usersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) CreateUser(ctx context.Context, user models.User) (int64, error) {
	const query = `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			phone,
			access_type
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var id int64

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.AccessType,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}
