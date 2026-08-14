package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Load loads configuration from YAML file and environment variables
// Environment variables take precedence over YAML file values
func Load(configFilePath string) (*Config, error) {
	k := koanf.New(".")

	// Load from YAML file if it exists
	if configFilePath != "" {
		if _, err := os.Stat(configFilePath); err == nil {
			if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
				return nil, fmt.Errorf("failed to load YAML config: %w", err)
			}
		}
	}

	// Load from environment variables (with underscore callback for nested keys)
	// Environment variables use underscores for nesting
	// e.g., SERVER_HOST, SERVER_PORT, DATABASE_PASSWORD
	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load env config: %w", err)
	}

	// Unmarshal into Config struct
	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// LoadWithDefaults loads configuration with default values
func LoadWithDefaults(configFilePath string) (*Config, error) {
	k := koanf.New(".")

	// Set default values
	defaults := map[string]interface{}{
		"app.name":          "user-service",
		"app.version":       "1.0.0",
		"app.env":           "development",
		"server.host":       "0.0.0.0",
		"server.port":       8080,
		"database.host":     "localhost",
		"database.port":     5432,
		"database.user":     "postgres",
		"database.password": "postgres",
		"database.database": "users_db",
		"database.sslmode":  "disable",
		"database.max_conn": 10,
		"database.min_conn": 2,
	}

	if err := k.Load(koanf.KeyMapProvider(defaults), koanf.Json()); err != nil {
		return nil, fmt.Errorf("failed to set defaults: %w", err)
	}

	// Load from YAML file if it exists
	if configFilePath != "" {
		if _, err := os.Stat(configFilePath); err == nil {
			if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
				return nil, fmt.Errorf("failed to load YAML config: %w", err)
			}
		}
	}

	// Load from environment variables (override file config)
	// Environment variables use underscores for nesting
	// e.g., SERVER_HOST, SERVER_PORT, DATABASE_PASSWORD
	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "_", "."))
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load env config: %w", err)
	}

	// Unmarshal into Config struct
	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
