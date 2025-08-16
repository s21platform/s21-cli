package devenv

import (
	"fmt"

	"github.com/BurntSushi/toml"
	userconfig "github.com/s21platform/s21-cli/internal/config"
)

// Config представляет структуру конфигурационного файла devenv.toml
type Config struct {
	Services ServicesConfig    `toml:"services"`
	Env      map[string]string `toml:"env,omitempty"`
	Creds    CredsConfig       `toml:"creds"`
	EnvMaps  []EnvMap          `toml:"env_maps,omitempty"`
}

// EnvMap описывает маппинг переменных окружения
type EnvMap struct {
	From string `toml:"from"` // Исходная переменная (например, POSTGRES_USER)
	To   string `toml:"to"`   // Целевая переменная (например, SERVICE_NAME_POSTGRES_USER)
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
	Topics  []Topic `toml:"topics,omitempty"`
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
		if c.Redpanda == nil {
			return false
		}
		if c.Redpanda.Topics == nil {
			c.Redpanda.Topics = []Topic{}
		}
		return c.Redpanda.Enabled
	case "redpanda_console":
		return c.RedpandaConsole != nil && c.RedpandaConsole.Enabled
	default:
		return false
	}
}

// LoadConfig загружает конфигурацию из TOML файла
func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("ошибка при чтении конфигурации: %v", err)
	}

	// Получаем пользовательскую конфигурацию
	userCfg, err := userconfig.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении конфигурации пользователя: %v", err)
	}

	// Добавляем nickname в переменные окружения
	if userCfg.Nickname != "" {
		if config.Env == nil {
			config.Env = make(map[string]string)
		}
		config.Env["ENV"] = userCfg.Nickname
	}

	// Применяем маппинг переменных
	if err := config.applyEnvMaps(); err != nil {
		return nil, fmt.Errorf("ошибка при маппинге переменных: %v", err)
	}

	return &config, nil
}

// applyEnvMaps применяет маппинг переменных окружения
func (c *Config) applyEnvMaps() error {
	if len(c.EnvMaps) == 0 {
		return nil
	}

	if c.Env == nil {
		c.Env = make(map[string]string)
	}

	// Создаем временную копию для новых значений
	newEnv := make(map[string]string)

	// Применяем каждый маппинг
	for _, mapping := range c.EnvMaps {
		if value, exists := c.Env[mapping.From]; exists {
			newEnv[mapping.To] = value
		}
	}

	// Добавляем новые значения в основную мапу
	for k, v := range newEnv {
		c.Env[k] = v
	}

	return nil
}
