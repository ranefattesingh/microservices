package config

type MainConfig struct {
	GinMode    string     `yaml:"gin_mode" mapstructure:"gin_mode" default:"test"`
	LogLevel   string     `yaml:"log_level" mapstructure:"log_level" default:"info"`
	HTTPConfig HTTPConfig `yaml:"http_config" mapstructure:"http_config"`
}

type HTTPConfig struct {
	Host string `yaml:"host" mapstructure:"host" default:"0.0.0.0"`
	Port int    `yaml:"port" mapstructure:"port" default:"8001"`
}
