package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ranefattesingh/ecommerce-platform/users/platform/database/psql"
	"github.com/ranefattesingh/ecommerce-platform/users/repository/models"
)

var ErrNoRowsUpdated = errors.New("no rows were updated")

type UsersRepository interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUserHavingEmailOrPhone(ctx context.Context, email, phone string) ([]models.User, error)
	GetUser(ctx context.Context, id int64) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
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
		return 0, fmt.Errorf("usersRepository.CreateUser: %w", err)
	}

	return id, nil
}

func (r *usersRepository) GetUsersHavingEmailOrPhone(ctx context.Context, email, phone string) ([]models.User, error) {
	const query = `
		SELECT id, first_name, last_name, email, phone, access_type, created_at, updated_at
		FROM users
		WHERE email = $1 OR phone = $2
	`
	rows, err := r.db.Pool.Query(ctx, query, email, phone)
	if err != nil {
		return nil, fmt.Errorf("usersRepository.GetUsersHavingEmailOrPhone: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("usersRepository.GetUsersHavingEmailOrPhone: %w", err)
	}
	return users, nil
}

func (r *usersRepository) GetUser(ctx context.Context, id int64) (models.User, error) {
	const query = `
		SELECT id, first_name, last_name, email, phone, access_type, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	rows, err := r.db.Pool.Query(ctx, query, id)
	if err != nil {
		return models.User{}, fmt.Errorf("usersRepository.GetUser: %w", err)
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return models.User{}, fmt.Errorf("usersRepository.GetUser: %w", err)
	}

	return user, nil
}

func (r *usersRepository) UpdateUser(ctx context.Context, user models.User) error {
	const query = `
		UPDATE users
		SET email = @newEmail, first_name=@newFirstName, last_name=@newLastName, phone=@newPhone ,updated_at = NOW()
		WHERE id = @userID
	`

	args := pgx.NamedArgs{
		"newFirstName": user.FirstName,
		"newLastName":  user.LastName,
		"newEmail":     user.Email,
		"newPhone":     user.Phone,
		"userID":       user.ID,
	}

	cmd, err := r.db.Pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("usersRepository.UpdateUser: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("usersRepository.UpdateUser: %w", ErrNoRowsUpdated)
	}

	return nil
}
