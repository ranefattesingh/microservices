package config

import (
	"strconv"
	"strings"
)

type MainConfig struct {
	GinMode        string         `yaml:"gin_mode" mapstructure:"gin_mode" default:"test"`
	LogLevel       string         `yaml:"log_level" mapstructure:"log_level" default:"info"`
	HTTPConfig     HTTPConfig     `yaml:"http" mapstructure:"http"`
	DatabaseConfig DatabaseConfig `yaml:"db" mapstructure:"db"`
}

type HTTPConfig struct {
	Host string `yaml:"host" mapstructure:"host" default:"0.0.0.0"`
	Port int    `yaml:"port" mapstructure:"port" default:"8001"`
}

type DatabaseConfig struct {
	Host     string             `yaml:"host" mapstructure:"host"`
	Port     string             `yaml:"port" mapstructure:"port"`
	Name     string             `yaml:"name" mapstructure:"name"`
	User     string             `yaml:"user" mapstructure:"user"`
	Password string             `yaml:"password" mapstructure:"password"`
	UseSSL   bool               `yaml:"use_ssl" mapstructure:"use_ssl" default:"false"`
	Pool     DatabasePoolConfig `yaml:"pool" mapstructure:"pool"`
}

type DatabasePoolConfig struct {
	ConnMaxLifetime       string `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime" default:"5m"`
	MaxConnIdleTime       string `yaml:"max_conn_idle_time" mapstructure:"max_conn_idle_time" default:"5m"`
	MaxConnLifeTimeJitter string `yaml:"max_conn_lifetime_jitter" mapstructure:"max_conn_lifetime_jitter" default:"30s"`
	MaxHealthCheckPeriod  string `yaml:"health_check_period" mapstructure:"health_check_period" default:"5s"`
	MaxConns              int    `yaml:"max_conns" mapstructure:"max_conns" default:"1000"`
	MinConns              int    `yaml:"min_conns" mapstructure:"min_conns" default:"1"`
}

func (dc DatabaseConfig) GetConnectionString() string {
	var builder strings.Builder

	builder.WriteString("postgresql://")
	builder.WriteString(dc.User)
	builder.WriteString(":" + dc.Password + "@")
	builder.WriteString(dc.Host)
	builder.WriteString(":" + dc.Port + "/")
	builder.WriteString(dc.Name)

	sslMode := "?sslmode=disable"
	if dc.UseSSL {
		sslMode = "?sslmode=enable"
	}

	builder.WriteString(sslMode)
	builder.WriteString("&pool_max_conns=" + strconv.Itoa(dc.Pool.MaxConns))
	builder.WriteString("&pool_min_conns=" + strconv.Itoa(dc.Pool.MinConns))
	builder.WriteString("&pool_max_conn_lifetime=" + dc.Pool.ConnMaxLifetime)
	builder.WriteString("&pool_max_conn_idle_time=" + dc.Pool.MaxConnIdleTime)
	builder.WriteString("&pool_health_check_period=" + dc.Pool.MaxHealthCheckPeriod)
	builder.WriteString("&pool_max_conn_lifetime_jitter=" + dc.Pool.MaxConnLifeTimeJitter)

	return builder.String()
}
