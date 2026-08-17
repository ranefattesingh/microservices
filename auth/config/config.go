package config

// ServerConfig holds server settings
type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

// Config represents the application configuration
type Config struct {
	Server ServerConfig `koanf:"server"`
}
