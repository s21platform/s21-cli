package devenv

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config представляет структуру конфигурационного файла devenv.toml
type Config struct {
	Services ServicesConfig    `toml:"services"`
	Env      map[string]string `toml:"env,omitempty"`
	Creds    CredsConfig       `toml:"creds"`
}

// CredsConfig конфигурация сервиса учетных данных
type CredsConfig struct {
	Enabled  bool     `toml:"enabled"`
	Endpoint string   `toml:"endpoint"`
	Services []string `toml:"services"`
}

// ServicesConfig содержит конфигурацию сервисов
type ServicesConfig struct {
	Postgres        *PostgresConfig        `toml:"postgres,omitempty"`
	Redis           *RedisConfig           `toml:"redis,omitempty"`
	Redpanda        *RedpandaConfig        `toml:"redpanda,omitempty"`
	RedpandaConsole *RedpandaConsoleConfig `toml:"redpanda_console,omitempty"`
}

// PostgresConfig конфигурация PostgreSQL
type PostgresConfig struct {
	Enabled bool `toml:"enabled"`
}

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Enabled bool `toml:"enabled"`
}

// RedpandaConfig конфигурация Redpanda
type RedpandaConfig struct {
	Enabled bool    `toml:"enabled"`
	Topics  []Topic `toml:"topics"`
}

// Topic конфигурация топика Kafka
type Topic struct {
	Name              string `toml:"name"`
	Partitions        int    `toml:"partitions"`
	ReplicationFactor int    `toml:"replication_factor"`
}

// RedpandaConsoleConfig конфигурация Redpanda Console
type RedpandaConsoleConfig struct {
	Enabled bool `toml:"enabled"`
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
	case "redpanda_console":
		return c.RedpandaConsole != nil && c.RedpandaConsole.Enabled
	default:
		return false
	}
}

// ValidateEnv проверяет наличие необходимых переменных окружения
func (c *Config) ValidateEnv() error {
	var missingVars []string

	// Проверяем переменные для PostgreSQL
	if c.Services.IsEnabled("postgres") {
		requiredVars := []string{
			"POSTGRES_USER",
			"POSTGRES_PASSWORD",
			"POSTGRES_DB",
			"POSTGRES_PORT",
		}
		for _, v := range requiredVars {
			if _, exists := c.Env[v]; !exists {
				missingVars = append(missingVars, v)
			}
		}
	}

	// Проверяем переменные для Redis
	if c.Services.IsEnabled("redis") {
		if _, exists := c.Env["REDIS_PORT"]; !exists {
			missingVars = append(missingVars, "REDIS_PORT")
		}
	}

	// Проверяем переменные для Redpanda
	if c.Services.IsEnabled("redpanda") {
		if _, exists := c.Env["KAFKA_PORT"]; !exists {
			missingVars = append(missingVars, "KAFKA_PORT")
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("не указаны обязательные переменные окружения в секции [env]:\n%s\n\nДобавьте их в секцию [env] вашего devenv.toml файла", strings.Join(missingVars, "\n"))
	}

	return nil
}

// LoadConfig загружает конфигурацию из TOML файла
func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("ошибка при чтении конфигурации: %v", err)
	}

	// Проверяем наличие необходимых переменных
	// if err := config.ValidateEnv(); err != nil {
	// 	return nil, err
	// }

	return &config, nil
}
