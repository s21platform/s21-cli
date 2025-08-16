# S21 CLI

CLI-утилита для разработчиков S21 Platform, предоставляющая инструменты для локальной разработки.

## Установка

```bash
go install github.com/s21platform/s21-cli@latest
```

## Авторизация

Для доступа к сервисам S21 необходимо авторизоваться:

```bash
s21 login
```

## Управление локальным окружением

### Конфигурация (devenv.toml)

```toml
# Конфигурация сервиса учетных данных
[creds]
enabled = true
endpoint = "localhost:50051"
# Список сервисов, для которых нужно получить учетные данные
services = [
    "notification",
    "user",
    "community"
]

# Переменные окружения
[env]
# Базовые настройки
APP_ENV = "local"
APP_DEBUG = "true"
APP_PORT = "8080"

# Настройки сервисов
POSTGRES_USER = "postgres"
POSTGRES_PASSWORD = "postgres"
POSTGRES_DB = "app"
POSTGRES_PORT = "5432"
KAFKA_PORT = "9092"
REDPANDA_CONSOLE_PORT = "18080"

# Инфраструктурные сервисы
[services.postgres]
enabled = true

[services.redis]
enabled = true

[services.redpanda]
enabled = true
# Список топиков для автоматического создания
topics = [
    { name = "notifications", partitions = 3, replication_factor = 1 },
    { name = "user-events", partitions = 3, replication_factor = 1 },
    { name = "community-events", partitions = 3, replication_factor = 1 }
]

[services.redpanda_console]
enabled = true  # Включает Redpanda Console UI на http://localhost:18080
```

### Запуск окружения

```bash
s21 devenv start
```

Команда:
1. Читает конфигурацию из `devenv.toml`
2. Получает учетные данные от сервиса credentials (если включено)
3. Генерирует `.env.s21-cli` файл
4. Создает `docker-compose.s21-cli.yml`
5. Запускает контейнеры через docker-compose
6. Создает топики в Redpanda (если указаны)

> **Примечание**: CLI использует отдельные файлы `.env.s21-cli` и `docker-compose.s21-cli.yml`, чтобы не конфликтовать с существующими файлами в проекте.

### Остановка окружения

```bash
s21 devenv finish
```

Команда:
1. Останавливает все контейнеры
2. Удаляет созданные контейнеры
3. Очищает временные файлы

## Доступные сервисы

После запуска окружения доступны следующие сервисы:

### PostgreSQL
- Порт: `${POSTGRES_PORT}` (по умолчанию 5432)
- Пользователь: `${POSTGRES_USER}`
- Пароль: `${POSTGRES_PASSWORD}`
- База данных: `${POSTGRES_DB}`

### Redis
- Порт: `${REDIS_PORT}` (по умолчанию 6379)

### Redpanda (Kafka API)
- Брокер: `localhost:${KAFKA_PORT}` (по умолчанию 9092)
- Admin API: `localhost:9644`
- Redpanda Console: `http://localhost:${REDPANDA_CONSOLE_PORT}` (по умолчанию 18080)
  - Пользователь: admin
  - Пароль: admin

## Команды

```bash
# Показать справку
s21 --help

# Показать версию
s21 version

# Авторизация
s21 login

# Управление окружением
s21 devenv --help       # Справка по командам окружения
s21 devenv start        # Запустить окружение
s21 devenv finish       # Остановить окружение
s21 devenv load        # Загрузить переменные окружения и запустить shell для разработки
```