package psql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, url string) (*Database, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	return &Database{Pool: pool}, nil
}
