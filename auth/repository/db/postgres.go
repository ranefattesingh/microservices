package db

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates and returns a Postgres connection pool using DATABASE_URL env var.
// Default: postgres://postgres:postgres@localhost:5432/authdb?sslmode=disable
func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	conn := os.Getenv("DATABASE_URL")
	if conn == "" {
		conn = "postgres://postgres:postgres@localhost:5432/authdb?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}
	// reasonable defaults
	cfg.MaxConns = 5
	cfg.MinConns = 0
	cfg.MaxConnLifetime = 30 * time.Minute

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx2, cfg)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
