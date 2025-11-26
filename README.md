# Сервис назначения ревьюеров для Pull Request'ов

Сервис для автоматического назначения ревьюеров на Pull Request

## Описание

Сервис предоставляет HTTP API для:
- Управления пользователями (создание, обновление, удаление, получение)
- Управления командами (создание, обновление, удаление, получение)
- Создания PR с автоматическим назначением ревьюеров (до 2 активных ревьюверов из команды автора)
- Переназначения ревьюверов (из команды заменяемого ревьювера)
- Объединения PR (идемпотентная операция)
- Получения списка PR, назначенных конкретному пользователю

## Используемый тек

- **Язык**: Go 1.24
- **База данных**: PostgreSQL 17
- **HTTP роутер**: Gorilla Mux
- **Миграции**: Goose
- **Контейнеризация**: Docker & Docker Compose

### Запуск с Docker Compose

1. Клонируем репозиторий:
```bash
git clone https://github.com/anguless/mr-reviewer
cd mr-reviewer
```

2. Запускаем сервис:
```bash
docker-compose up
```

Сервис будет доступен на `http://localhost:8080`

*При запуске сервиса Миграции применяются автоматически.*

## Описание ручек

### Health Check
- `GET /health` - Проверка готовности сервиса

### Users
- `POST /api/v1/users` - Создать пользователя
- `GET /api/v1/users/{user_id}` - Получить пользователя
- `PUT /api/v1/users/{user_id}` - Обновить пользователя
- `DELETE /api/v1/users/{user_id}` - Удалить пользователя
- `GET /api/v1/users/{user_id}/pull-requests` - Получить список PR, назначенных пользователю

### Teams
- `POST /api/v1/teams` - Создать команду
- `GET /api/v1/teams/{team_id}` - Получить команду
- `PUT /api/v1/teams/{team_id}` - Обновить команду
- `DELETE /api/v1/teams/{team_id}` - Удалить команду

### Pull Requests
- `POST /api/v1/pull-requests` - Создать PR (автоматически назначаются ревьюверы)
- `GET /api/v1/pull-requests` - Получить все PR
- `GET /api/v1/pull-requests/{pull_request_id}` - Получить PR по ID
- `POST /api/v1/pull-requests/{pull_request_id}/reassign` - Переназначить ревьювера
- `POST /api/v1/pull-requests/{pull_request_id}/merge` - Объединить PR

## Примеры использования

### Создание команды
```bash
curl -X POST http://localhost:8080/api/v1/team/add \
  -H "Content-Type: application/json" \
  -d '
  {
  "team_name": "team",
  "members": [
    {
      "username": "Ivan Ivanov",
      "is_active": true
    }
  ]
}'
```

### Создание пользователя
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "team_id": "<team_id>",
    "is_active": true
  }'
```

### Создание PR
```bash
curl -X POST http://localhost:8080/api/v1/pull-requests \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_name": "Add new feature",
    "author_id": "<user_id>"
  }'
```

При создании PR автоматически назначаются до 2 активных ревьюверов из команды автора (исключая самого автора).

### Переназначение ревьювера
```bash
curl -X POST http://localhost:8080/api/v1/pull-requests/<pr_id>/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "reviewer_id": "<reviewer_id_to_replace>"
  }'
```

Новый ревьювер выбирается случайным образом из команды заменяемого ревьювера.

### Объединение PR
```bash
curl -X POST http://localhost:8080/api/v1/pull-requests/<pr_id>/merge
```

Операция идемпотентна - повторный вызов не приводит к ошибке.

## Taskfile
В проекте используется taskfile (в качестве альтернативы makefile)

Taskfile устанавливается следующей командой:
```bash
go install github.com/go-task/task/v3/cmd/task@latest
```
- `task install-formatters` - Устанавливает форматтеры gci и gofumpt в ./bin
- `task format` - Форматирует весь проект gofumpt + gci, исключая mocks
- `task install-golangci-lint` - Устанавливает golangci-lint в каталог bin
- `task lint` - Запускает golangci-lint + ```task format``` для всех модулей

### ✅ Конфигурация линтера

Настроен `golangci-lint` с конфигурацией в `.golangci.yml`, а также форматтеры gofumpt и gci

- Включены основные линтеры: errcheck, govet, golint, gosec, staticcheck
- Настроены правила для тестовых файлов
- Оптимизированы настройки для проекта

**Запуск линтера и форматтеров:**
```bash
task lint
```
