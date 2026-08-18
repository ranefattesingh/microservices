package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ranefattesingh/microservices/auth/repository/db/models"
	"github.com/ranefattesingh/microservices/pkg/pgx/pool"
	"github.com/redis/go-redis/v9"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

var _ AuthRepository = (*authRepository)(nil)

type AuthRepository interface {
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	SaveRefreshToken(ctx context.Context, token, email string, ttl time.Duration) error
	CreateAccount(ctx context.Context, account *models.Account) (int64, error)
	GetEmailByRefreshToken(ctx context.Context, token string) (string, error)
}

type authRepository struct {
	db    pool.Database
	cache *redis.Client
}

func NewAuthRepository(db pool.Database, cache *redis.Client) AuthRepository {
	return &authRepository{
		db:    db,
		cache: cache,
	}
}

func (r *authRepository) CreateAccount(ctx context.Context, account *models.Account) (int64, error) {
	// insert into Postgres
	var id int64
	err := r.db.Pool().QueryRow(ctx, "INSERT INTO accounts (email, password, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id", account.Email, account.Password).Scan(&id)
	if err != nil {
		// detect unique violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, ErrEmailAlreadyExists
			}
		}
		return 0, err
	}

	// cache account in redis for fast reads
	if r.cache != nil {
		key := fmt.Sprintf("user:%s", account.Email)
		_ = r.cache.HSet(ctx, key, map[string]interface{}{
			"id":       strconv.FormatInt(id, 10),
			"email":    account.Email,
			"password": account.Password,
		}).Err()
		// persist key
		_ = r.cache.Persist(ctx, key)
	}

	return id, nil
}

func (r *authRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	// try cache first
	if r.cache != nil {
		key := fmt.Sprintf("user:%s", email)
		vals, err := r.cache.HGetAll(ctx, key).Result()
		if err == nil && len(vals) > 0 {
			idStr := vals["id"]
			idInt, _ := strconv.ParseInt(idStr, 10, 64)
			return &models.Account{
				ID:       idInt,
				Email:    email,
				Password: vals["password"],
			}, nil
		}
	}

	// fallback to DB
	account := &models.Account{}
	err := r.db.Pool().QueryRow(ctx, "SELECT id, email, password, created_at, updated_at FROM accounts WHERE email=$1", email).Scan(&account.ID, &account.Email, &account.Password, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// populate cache for future
	if r.cache != nil {
		key := fmt.Sprintf("user:%s", email)
		_ = r.cache.HSet(ctx, key, map[string]interface{}{
			"id":       strconv.FormatInt(account.ID, 10),
			"email":    account.Email,
			"password": account.Password,
		}).Err()
		_ = r.cache.Persist(ctx, key)
	}

	return account, nil
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, token, email string, ttl time.Duration) error {
	if r.cache == nil {
		return errors.New("cache not configured")
	}
	key := fmt.Sprintf("refresh:%s", token)
	return r.cache.Set(ctx, key, email, ttl).Err()
}

func (r *authRepository) GetEmailByRefreshToken(ctx context.Context, token string) (string, error) {
	if r.cache == nil {
		return "", errors.New("cache not configured")
	}
	key := fmt.Sprintf("refresh:%s", token)
	res, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("invalid refresh token")
		}
		return "", err
	}
	return res, nil
}
