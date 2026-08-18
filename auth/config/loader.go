package config

import (
	"fmt"

	"github.com/ranefattesingh/microservices/pkg/config"
)

func LoadConfig() (*Config, error) {
	var conf *Config

	loader := config.NewLoader("config.yaml", "AUTH_SERVICE_")

	err := loader.Load(conf)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return conf, nil
}
