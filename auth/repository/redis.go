package repository

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and returns a go-redis client configured from environment
func NewRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	db := 0
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Basic ping to verify connectivity (non-fatal here)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = client.Ping(ctx).Err()

	return client
}
