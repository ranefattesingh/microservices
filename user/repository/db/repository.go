package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ranefattesingh/microservices/user/repository/db/models"
	"github.com/ranefattesingh/microservices/user/repository/db/pool"
)

var _ UserRepository = (*userRepository)(nil)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id int64, user *models.User) error
	GetAll(ctx context.Context, page, limit int) ([]*models.User, int, error)
}

type userRepository struct {
	db pool.Database
}

func NewUserRepository(db pool.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (int64, error) {
	query := `
		INSERT INTO users (first_name, last_name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`

	var id int64
	err := r.db.Pool().QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.Pool().QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.Pool().QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, id int64, user *models.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, password = $4, updated_at = NOW()
		WHERE id = $5
	`

	_, err := r.db.Pool().Exec(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		id,
	)

	return err
}

func (r *userRepository) GetAll(ctx context.Context, page, limit int) ([]*models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	countQuery := `SELECT COUNT(*) FROM users`
	var total int
	if err := r.db.Pool().QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * limit
	rows, err := r.db.Pool().Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.User])
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
