package config

// DefaultLoader returns a loader with sensible defaults for local development.
func DefaultLoader() Loader {
	return Loader{}
}

type Loader struct{}

// Load returns default configuration. Keep it simple to match project style.
func (l Loader) Load() (Config, error) {
	return Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8081,
		},
	}, nil
}
