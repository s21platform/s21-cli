package devenv

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const composeTemplate = `version: '3.8'

services:
{{- if .Services.IsEnabled "postgres" }}
  postgres:
    image: postgres:{{ .Services.Postgres.Version }}
    container_name: dev-postgres
    environment:
      POSTGRES_USER: {{ .Services.Postgres.User }}
      POSTGRES_PASSWORD: {{ .Services.Postgres.Password }}
      POSTGRES_DB: {{ .Services.Postgres.Database }}
    ports:
      - "{{ .Services.Postgres.Port }}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - dev_network
{{- end }}

{{- if .Services.IsEnabled "redis" }}
  redis:
    image: redis:{{ .Services.Redis.Version }}
    container_name: dev-redis
    ports:
      - "{{ .Services.Redis.Port }}:6379"
    volumes:
      - redis_data:/data
    networks:
      - dev_network
{{- end }}

{{- if .Services.IsEnabled "redpanda" }}
  redpanda:
    image: redpandadata/redpanda:{{ .Services.Redpanda.Version }}
    container_name: dev-redpanda
    ports:
      - "{{ .Services.Redpanda.Port }}:9092"
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

// GenerateCompose генерирует docker-compose.yml файл на основе конфигурации
func GenerateCompose(config *Config, outputPath string) error {
	tmpl, err := template.New("compose").Parse(composeTemplate)
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
	if err := tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("ошибка при генерации файла: %v", err)
	}

	return nil
}
