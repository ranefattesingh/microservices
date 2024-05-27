package config

import (
	"strconv"
	"strings"
)

type MainConfig struct {
	GinMode        string         `default:"test"      mapstructure:"gin_mode"  yaml:"gin_mode"`
	LogLevel       string         `default:"info"      mapstructure:"log_level" yaml:"log_level"`
	HTTPConfig     HTTPConfig     `mapstructure:"http" yaml:"http"`
	DatabaseConfig DatabaseConfig `mapstructure:"db"   yaml:"db"`
}

type HTTPConfig struct {
	Host string `default:"0.0.0.0" mapstructure:"host" yaml:"host"`
	Port int    `default:"8001"    mapstructure:"port" yaml:"port"`
}

type DatabaseConfig struct {
	Host     string             `mapstructure:"host"     yaml:"host"`
	Port     string             `mapstructure:"port"     yaml:"port"`
	Name     string             `mapstructure:"name"     yaml:"name"`
	User     string             `mapstructure:"user"     yaml:"user"`
	Password string             `mapstructure:"password" yaml:"password"`
	UseSSL   bool               `default:"false"         mapstructure:"use_ssl" yaml:"use_ssl"`
	Pool     DatabasePoolConfig `mapstructure:"pool"     yaml:"pool"`
}

type DatabasePoolConfig struct {
	ConnMaxLifetime       string `default:"5m"   mapstructure:"conn_max_lifetime"        yaml:"conn_max_lifetime"`
	MaxConnIdleTime       string `default:"5m"   mapstructure:"max_conn_idle_time"       yaml:"max_conn_idle_time"`
	MaxConnLifeTimeJitter string `default:"30s"  mapstructure:"max_conn_lifetime_jitter" yaml:"max_conn_lifetime_jitter"`
	MaxHealthCheckPeriod  string `default:"5s"   mapstructure:"health_check_period"      yaml:"health_check_period"`
	MaxConns              int    `default:"1000" mapstructure:"max_conns"                yaml:"max_conns"`
	MinConns              int    `default:"1"    mapstructure:"min_conns"                yaml:"min_conns"`
}

//nolint:revive
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

	return builder.String()
}

//nolint:revive
func (dc DatabaseConfig) GetConnectionStringWithOptions() string {
	var builder strings.Builder

	builder.WriteString(dc.GetConnectionString())

	builder.WriteString("&pool_max_conns=" + strconv.Itoa(dc.Pool.MaxConns))
	builder.WriteString("&pool_min_conns=" + strconv.Itoa(dc.Pool.MinConns))
	builder.WriteString("&pool_max_conn_lifetime=" + dc.Pool.ConnMaxLifetime)
	builder.WriteString("&pool_max_conn_idle_time=" + dc.Pool.MaxConnIdleTime)
	builder.WriteString("&pool_health_check_period=" + dc.Pool.MaxHealthCheckPeriod)
	builder.WriteString("&pool_max_conn_lifetime_jitter=" + dc.Pool.MaxConnLifeTimeJitter)

	return builder.String()
}
