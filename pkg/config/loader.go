package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type loader struct {
	path string
}

func DefaultLoader() *loader {
	return &loader{
		path: "config.yaml",
	}
}

func (l *loader) Load() (*Config, error) {
	k := koanf.New(".")

	loaded, err := l.loadFile(k)
	if err != nil {
		return nil, err
	}

	if !loaded {
		if err := l.loadEnv(k); err != nil {
			return nil, err
		}
	}

	cfg := new(Config)

	if err := k.Unmarshal("", cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func (l *loader) loadFile(k *koanf.Koanf) (bool, error) {
	if l.path == "" {
		return false, nil
	}

	if _, err := os.Stat(l.path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("stat config file: %w", err)
	}

	if err := k.Load(
		file.Provider(l.path),
		yaml.Parser(),
	); err != nil {
		return false, fmt.Errorf("load config file: %w", err)
	}

	return true, nil
}

func (l *loader) loadEnv(k *koanf.Koanf) error {
	const prefix = ""

	if err := k.Load(
		env.Provider(
			prefix,
			".",
			func(s string) string {
				s = strings.TrimPrefix(s, prefix)

				return strings.ToLower(
					strings.ReplaceAll(s, "__", "."),
				)
			},
		),
		nil,
	); err != nil {
		return fmt.Errorf("load environment config: %w", err)
	}

	return nil
}
