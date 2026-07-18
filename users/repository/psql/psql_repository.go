package psql

import (
	"context"
	"fmt"

	"github.com/ranefattesingh/ecommerce-platform/users/platform/database/psql"
	"github.com/ranefattesingh/ecommerce-platform/users/repository/models"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUserHavingEmailOrPhone(ctx context.Context, email, phone string) (models.User, error)
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

func (r *usersRepository) GetUserHavingEmailOrPhone(ctx context.Context, email, phone string) (models.User, error) {
	const query = `
		SELECT id, email, phone
		FROM users
		WHERE email = $1 OR phone = $2
	`

	var user models.User

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		email,
		phone,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
	)

	if err != nil {
		return user, fmt.Errorf("read user %w", err)
	}

	return user, nil
}
