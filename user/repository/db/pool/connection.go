package pool

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	Pool() *pgxpool.Pool
}

type connection struct {
	pool *pgxpool.Pool
}

func NewConnectionPool(conStr string) (*connection, error) {
	pool, err := pgxpool.New(context.Background(), conStr)

	if err != nil {
		return nil, nil
	}

	c := &connection{
		pool: pool,
	}

	return c, nil
}

func (c *connection) Close() error {
	return c.Close()
}

func (c *connection) Pool() *pgxpool.Pool {
	return c.pool
}
