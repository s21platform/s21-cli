package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Token string `json:"token"`
}

// GetConfigDir возвращает путь к директории конфигурации
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("не удалось получить домашнюю директорию: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "s21-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию конфигурации: %v", err)
	}

	return configDir, nil
}

// GetConfigPath возвращает путь к файлу конфигурации
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "creds.json"), nil
}

// LoadConfig загружает конфигурацию
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл конфигурации: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("не удалось разобрать файл конфигурации: %v", err)
	}

	return &config, nil
}

// SaveConfig сохраняет конфигурацию
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать конфигурацию: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("не удалось сохранить файл конфигурации: %v", err)
	}

	return nil
}
