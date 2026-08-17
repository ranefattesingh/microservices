package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

// UserRepository provides persistence using Postgres for durable storage
// and Redis for caching and refresh tokens.
// Keying strategy in cache:
// - user:{email} -> hash fields: id, email, password
// - refresh:{token} -> string email (with TTL)
type UserRepository struct {
	db    *pgxpool.Pool
	cache *redis.Client
}

func NewUserRepository(db *pgxpool.Pool, cache *redis.Client) *UserRepository {
	return &UserRepository{db: db, cache: cache}
}

func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash string) (int64, error) {
	// insert into Postgres
	var id int64
	err := r.db.QueryRow(ctx, "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", email, passwordHash).Scan(&id)
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

	// cache user in redis for fast reads
	if r.cache != nil {
		key := fmt.Sprintf("user:%s", email)
		_ = r.cache.HSet(ctx, key, map[string]interface{}{
			"id":       strconv.FormatInt(id, 10),
			"email":    email,
			"password": passwordHash,
		}).Err()
		// persist key
		_ = r.cache.Persist(ctx, key)
	}

	return id, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (id int64, storedHash string, err error) {
	// try cache first
	if r.cache != nil {
		key := fmt.Sprintf("user:%s", email)
		vals, err := r.cache.HGetAll(ctx, key).Result()
		if err == nil && len(vals) > 0 {
			idStr := vals["id"]
			idInt, _ := strconv.ParseInt(idStr, 10, 64)
			return idInt, vals["password"], nil
		}
	}

	// fallback to DB
	var idInt int64
	var hash string
	err = r.db.QueryRow(ctx, "SELECT id, password FROM users WHERE email=$1", email).Scan(&idInt, &hash)
	if err != nil {
		return 0, "", ErrUserNotFound
	}

	// populate cache for future
	if r.cache != nil {
		key := fmt.Sprintf("user:%s", email)
		_ = r.cache.HSet(ctx, key, map[string]interface{}{
			"id":       strconv.FormatInt(idInt, 10),
			"email":    email,
			"password": hash,
		}).Err()
		_ = r.cache.Persist(ctx, key)
	}

	return idInt, hash, nil
}

func (r *UserRepository) StoreRefreshToken(ctx context.Context, token, email string, ttl time.Duration) error {
	if r.cache == nil {
		return errors.New("cache not configured")
	}
	key := fmt.Sprintf("refresh:%s", token)
	return r.cache.Set(ctx, key, email, ttl).Err()
}

func (r *UserRepository) GetEmailByRefreshToken(ctx context.Context, token string) (string, error) {
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
