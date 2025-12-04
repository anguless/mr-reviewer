# Сервис назначения ревьюеров для Pull Request'ов

***Тестовое задание для backend-стажировки в Авито, осень 2025***

Сервис "PR Reviewer Assignment Service" предназначен для автоматизации процесса назначения ревьюеров к Pull Requests (PR) в командах разработки.
Он позволяет создавать команды с участниками, отмечать активность пользователей, создавать и отслеживать статус пулл-реквестов,
а также управлять назначением и перераспределением ревьюверов в пределах команды.

## Описание

Основные функции данного сервиса:
- Управление командами и их участниками (создание команд, получение информации о командах).
- Управление пользователями с возможностью установки статуса активности.
- Управление пулл-реквестами: создание новых PR, слияние (merge) PR, получение списка PR для обзора, перераспределение ревьюверов.
- Предоставление информации о назначенных PR для конкретного пользователя для удобства управления ревьюванием кода.

*Проект построен с использованием принципов чистой архитектуры*

## Используемый cтек

- **Язык**: Go 1.24.7
- **База данных**: PostgreSQL 17
- **HTTP роутер**: Chi v5.2.3
- **Миграции**: Goose v3.26
- **Кодогенератор**: В проекте используется OpenAPI-кодогенератор [Ogen](https://github.com/ogen-go/ogen).
  Документация, согласно которой генерировались файлы, находится по пути: */shared/api/mr/v1/mr.openapi.yml*.
  Все сгенерированные с помощью Ogen файлы находятся в папке: */pkg/openapi/mr/v1*
- **Контейнеризация**: Docker & Docker Compose

### Запуск с Docker Compose

1. Клонируем репозиторий:
```bash
git clone https://github.com/anguless/mr-reviewer
cd mr-reviewer
```

2. Запускаем сервис:
```bash
docker-compose up -d
```

Сервис будет доступен на `http://localhost:8080`

*При запуске сервиса Миграции применяются автоматически.*

## Описание ручек

### Teams
- `/team/add` - `POST` - Создание команды с именем и списком участников.

***Пример***
```bash
curl -X 'POST' \
  'http://localhost:8080/team/add' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "team_name": "payments",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": true
      }
    ]
  }'
```
- `/team/get` - `GET`. Получение информации о команде по параметру ?team_name=.
  
***Пример***
```bash
curl -X 'GET' \
  'http://localhost:8080/team/get?team_name=backend' \
  -H 'accept: application/json'
```

### Users
- `/users/setIsActive` - `POST`. Установка активности пользователя (is_active true/false) по user_id.
  
***Пример***
```bash
curl -X 'POST' \
  'http://localhost:8080/users/setIsActive' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": "u2",
    "is_active": false
  }'

```
- `/users/getReview` - `GET`. Получение списка PR для ревью по параметру ?userid=.

***Пример***
```bash
curl -X 'GET' \
  'http://localhost:8080/users/getReview?user_id=u2' \
  -H 'accept: application/json'

```

### Pull Requests
- `/pullRequest/create` - `POST`. Создание Pull Request'а с ID, названием и автором.
    При создании PR автоматически назначаются до 2 активных ревьюверов из команды автора (исключая самого автора).

***Пример***
```bash
curl -X 'POST' \
  'http://localhost:8080/pullRequest/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search",
    "author_id": "u1"
  }'
```

- `/pullRequest/merge` - `POST`. Слияние (merge) Pull Request'а по ID, обновляет статус на MERGED.
  
***Пример***
```bash
curl -X 'POST' \
  'http://localhost:8080/pullRequest/merge' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "pull_request_id": "pr-1001"
  }'
```

- `/pullRequest/reassign` - `POST`. Замена ревьювера (old_user_id) на другого активного в команде. Новый ревьювер выбирается случайным образом из команды заменяемого ревьювера.
    Операция идемпотентна - повторный вызов не приводит к ошибке.

***Пример***
```bash
curl -X 'POST' \
  'http://localhost:8080/pullRequest/reassign' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "pull_request_id": "pr-1001",
    "old_user_id": "u2"
  }'
```

## Taskfile
В проекте используется ***Taskfile*** (в качестве альтернативы ***Makefile***)

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

### Примечание

Файл **.env** добавлен в репозиторий намеренно, а не по случайности или невнимательности :)
Сделано это в соответствие с условиями и правилами выполнения данного тестового задания, а также для удобства проверяющих.