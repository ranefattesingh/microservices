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
	path   string
	prefix string
}

func NewLoader(path, prefix string) *loader {
	return &loader{
		path:   path,
		prefix: prefix,
	}
}

func (l *loader) Load(c any) error {
	k := koanf.New(".")

	loaded, err := l.loadFile(k)
	if err != nil {
		return err
	}

	if !loaded {
		if err := l.loadEnv(k); err != nil {
			return err
		}
	}

	if err := k.Unmarshal("", c); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
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
	if err := k.Load(
		env.Provider(
			l.prefix,
			".",
			func(s string) string {
				s = strings.TrimPrefix(s, l.prefix)

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
