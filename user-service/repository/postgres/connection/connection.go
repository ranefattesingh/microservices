package connection

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection interface {
	Pool() *pgxpool.Pool
	TestConnection() error
}

type connection struct {
	pool *pgxpool.Pool
}

func NewConnectionPool(connStr string) (*connection, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	connection := &connection{
		pool: pool,
	}

	return connection, nil
}

func (c *connection) Pool() *pgxpool.Pool {
	return c.pool
}

func (c *connection) TestConnection() error {
	return c.pool.Ping(context.Background())
}
