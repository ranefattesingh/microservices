package config

import (
	"fmt"
	"net/url"
)

// Config represents the application configuration
type Config struct {
	App      AppConfig      `koanf:"app"`
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
}

// AppConfig holds application-level settings
type AppConfig struct {
	Name    string `koanf:"name"`
	Version string `koanf:"version"`
	Env     string `koanf:"env"`
}

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

func (dc DatabaseConfig) EncodedConnectionString() string {
	u := &url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%d", dc.Host, dc.Port),
		Path:   dc.Database,
	}

	u.User = url.UserPassword(dc.User, dc.Password)

	q := u.Query()
	if dc.SSLMode != "" {
		q.Set("sslmode", dc.SSLMode)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
