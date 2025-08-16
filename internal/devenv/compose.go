package devenv

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// generateRandomSuffix генерирует случайный суффикс для имен контейнеров
func generateRandomSuffix() string {
	bytes := make([]byte, 4) // 8 символов в hex
	if _, err := rand.Read(bytes); err != nil {
		return "fallback"
	}
	return hex.EncodeToString(bytes)
}

const composeTemplate = `# yaml-language-server: $schema=https://raw.githubusercontent.com/compose-spec/compose-spec/master/schema/compose-spec.json
version: '3.8'
# Автоматически сгенерированный файл, не редактируйте его вручную!
# Используйте devenv.toml для настройки окружения.

services:
{{- if .Services.IsEnabled "postgres" }}
  postgres:
    image: postgres:15
    container_name: {{ containerName "postgres" }}
    env_file:
      - .env.s21-cli
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "${POSTGRES_PORT}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - dev_network
{{- end }}

{{- if .Services.IsEnabled "redis" }}
  redis:
    image: redis:alpine
    container_name: {{ containerName "redis" }}
    env_file:
      - .env.s21-cli
    ports:
      - "${REDIS_PORT}:6379"
    volumes:
      - redis_data:/data
    networks:
      - dev_network
{{- end }}

{{- if .Services.IsEnabled "redpanda" }}
  redpanda:
    image: redpandadata/redpanda:latest
    container_name: {{ containerName "redpanda" }}
    env_file:
      - .env.s21-cli
    ports:
      - "${KAFKA_PORT}:9092"
      - "9644:9644"
    volumes:
      - redpanda_data:/var/lib/redpanda/data
    networks:
      - dev_network
    command:
      - redpanda
      - start
      - --overprovisioned
      - --smp
      - "1"
      - --memory
      - "1G"
      - --reserve-memory
      - "0M"
      - --node-id
      - "0"
      - --check=false
      - --kafka-addr
      - PLAINTEXT://0.0.0.0:9092
      - --advertise-kafka-addr
      - PLAINTEXT://redpanda:9092
    healthcheck:
      test: ["CMD", "rpk", "cluster", "health"]
      interval: 10s
      timeout: 5s
      retries: 5
{{- end }}

{{- if and (.Services.IsEnabled "redpanda") .Services.Redpanda.Topics }}
  init-topics:
    image: redpandadata/redpanda:latest
    container_name: {{ containerName "init-topics" }}
    depends_on:
      redpanda:
        condition: service_healthy
        required: true
    networks:
      - dev_network
    entrypoint: ["/bin/sh"]
    command:
      - -c
      - |
        {{- range .Services.Redpanda.Topics }}
        echo "Creating topic {{.Name}}..."
        rpk topic create {{.Name}} --brokers redpanda:9092 --partitions {{.Partitions}} --replicas {{.ReplicationFactor}} || echo "Topic {{.Name}} might already exist"
        {{- end }}
        echo "All topics have been created!"
{{- end }}

{{- if .Services.IsEnabled "redpanda_console" }}
  redpanda-console:
    image: docker.redpanda.com/redpandadata/console:latest
    container_name: {{ containerName "redpanda-console" }}
    env_file:
      - .env.s21-cli
    ports:
      - "${REDPANDA_CONSOLE_PORT}:8080"
    environment:
      KAFKA_BROKERS: redpanda:9092
      CONSOLE_BASIC_AUTH_USERNAME: admin
      CONSOLE_BASIC_AUTH_PASSWORD: admin
    depends_on:
      redpanda:
        condition: service_healthy
        required: true
    networks:
      - dev_network
{{- end }}

volumes:
{{- if .Services.IsEnabled "postgres" }}
  postgres_data:
{{- end }}
{{- if .Services.IsEnabled "redis" }}
  redis_data:
{{- end }}
{{- if .Services.IsEnabled "redpanda" }}
  redpanda_data:
{{- end }}

networks:
  dev_network:
    driver: bridge`

// ComposeData содержит данные для генерации docker-compose файла
type ComposeData struct {
	*Config
	ContainerNames map[string]string // Маппинг сервис -> имя контейнера
}

// GenerateCompose генерирует docker-compose.yml файл на основе конфигурации
func GenerateCompose(config *Config, outputPath string) error {
	// Генерируем уникальный суффикс для этого запуска
	// Генерируем уникальные имена для контейнеров
	suffix := generateRandomSuffix()
	containerNames := make(map[string]string)
	containerNames["postgres"] = fmt.Sprintf("dev-postgres-%s", suffix)
	containerNames["redis"] = fmt.Sprintf("dev-redis-%s", suffix)
	containerNames["redpanda"] = fmt.Sprintf("dev-redpanda-%s", suffix)
	containerNames["init-topics"] = fmt.Sprintf("dev-init-topics-%s", suffix)
	containerNames["redpanda-console"] = fmt.Sprintf("dev-redpanda-console-%s", suffix)

	data := ComposeData{
		Config:         config,
		ContainerNames: containerNames,
	}
	funcMap := template.FuncMap{
		"containerName": func(service string) string {
			return data.ContainerNames[service]
		},
	}
	tmpl, err := template.New("compose").Funcs(funcMap).Parse(composeTemplate)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге шаблона: %v", err)
	}

	// Создаем директорию, если она не существует
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("ошибка при создании директории: %v", err)
	}

	// Создаем файл
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %v", err)
	}
	defer file.Close()

	// Генерируем содержимое
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("ошибка при генерации файла: %v", err)
	}

	return nil
}
