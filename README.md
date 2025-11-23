# PR_service

# Инструкция по запуску проекта

## Требования
- Docker
- Docker Compose

## Сборка и запуск

### 1. Собрать и запустить сервисы
В корне проекта выполнить:

```docker-compose up --build```

Будут подняты два контейнера:
- `go_service` — приложение
- `postgres` — база данных PostgreSQL

Приложение будет доступно на порту `8080`.

### 2. Переменные окружения
Приложение использует переменную:

```DATABASE_URL=postgres://postgres:postgres@postgres:5432/appdb?sslmode=disable```

Она задаётся в `docker-compose.yaml`.

### 3. Миграции
Контейнер содержит инструмент `goose` и каталог `migrations`.

После запуска контейнера миграции применяются входным скриптом `entrypoint.sh`.
При необходимости можно выполнить миграции вручную:

```docker exec -it go_service goose -dir /app/migrations postgres "$DATABASE_URL" up```

### 4. Остановка сервисов

```docker-compose down```

### 5. Хранение данных
Данные PostgreSQL сохраняются в volume:

```pgdata:/var/lib/postgresql/data```

Это позволяет сохранять состояние между перезапусками контейнеров.


В ходе решения задачи возникла проблема с миграцией БД, с выбором автогенератора кода для указанного апи