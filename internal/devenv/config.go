package devenv

import (
	"github.com/BurntSushi/toml"
)

// Config представляет структуру конфигурационного файла devenv.toml
type Config struct {
	Services ServicesConfig `toml:"services"`
}

// ServicesConfig содержит конфигурацию сервисов
type ServicesConfig struct {
	Postgres *PostgresConfig `toml:"postgres,omitempty"`
	Redis    *RedisConfig    `toml:"redis,omitempty"`
	Redpanda *RedpandaConfig `toml:"redpanda,omitempty"`
}

// IsEnabled проверяет, включен ли сервис
func (c *ServicesConfig) IsEnabled(service string) bool {
	switch service {
	case "postgres":
		return c.Postgres != nil && c.Postgres.Enabled
	case "redis":
		return c.Redis != nil && c.Redis.Enabled
	case "redpanda":
		return c.Redpanda != nil && c.Redpanda.Enabled
	default:
		return false
	}
}

// PostgresConfig конфигурация PostgreSQL
type PostgresConfig struct {
	Enabled  bool   `toml:"enabled"`
	Version  string `toml:"version"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Enabled bool   `toml:"enabled"`
	Version string `toml:"version"`
	Port    int    `toml:"port"`
}

// RedpandaConfig конфигурация Redpanda
type RedpandaConfig struct {
	Enabled bool   `toml:"enabled"`
	Version string `toml:"version"`
	Port    int    `toml:"port"`
}

// LoadConfig загружает конфигурацию из TOML файла
func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
