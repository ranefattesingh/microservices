package config

import (
	"fmt"
	"time"
)

// ServerConfig holds server settings
type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
	SSLMode  string `koanf:"sslmode"`
	MaxConn  int    `koanf:"max_conn"`
	MinConn  int    `koanf:"min_conn"`
}

// JWTConfig holds JWT settings
type JWTConfig struct {
	Secret          string        `koanf:"secret"`
	AccessTokenTTL  time.Duration `koanf:"access_token_ttl"`
	RefreshTokenTTL time.Duration `koanf:"refresh_token_ttl"`
}

// EncodedConnectionString returns a PostgreSQL connection string
func (d DatabaseConfig) EncodedConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", d.User, d.Password, d.Host, d.Port, d.Database, d.SSLMode)
}

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	JWT      JWTConfig      `koanf:"jwt"`
}
