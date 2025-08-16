package devenv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateEnvFile создает .env файл на основе конфигурации
func GenerateEnvFile(config *Config, outputPath string) error {
	// Генерируем уникальный суффикс для этого запуска
	suffix := generateRandomSuffix()
	// Создаем базовые переменные окружения из конфигурации сервисов
	envVars := make(map[string]string)

	// Добавляем базовые переменные для хостов
	if config.Services.IsEnabled("postgres") {
		envVars["POSTGRES_HOST"] = "localhost"
	}
	if config.Services.IsEnabled("redis") {
		envVars["REDIS_HOST"] = "localhost"
	}
	if config.Services.IsEnabled("redpanda") {
		envVars["KAFKA_HOST"] = "localhost"
	}

	// Добавляем суффикс для уникальных имен контейнеров
	envVars["DEVENV_SUFFIX"] = suffix

	// Добавляем пользовательские переменные (они могут перезаписать автоматически созданные)
	for key, value := range config.Env {
		envVars[key] = value
	}

	// Создаем директорию, если она не существует
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("ошибка при создании директории: %v", err)
	}

	// Создаем .env файл
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка при создании .env файла: %v", err)
	}
	defer file.Close()

	// Записываем переменные в файл
	for key, value := range envVars {
		if _, err := fmt.Fprintf(file, "%s=%s\n", key, value); err != nil {
			return fmt.Errorf("ошибка при записи в .env файл: %v", err)
		}
	}

	return nil
}

// LoadEnvFile загружает переменные окружения из .env файла
func LoadEnvFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ошибка при чтении .env файла: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("ошибка при установке переменной окружения %s: %v", key, err)
		}
	}

	return nil
}
