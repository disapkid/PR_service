# Стейдж сборки
FROM golang:1.24 AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Сборка бинаря
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# Финальный стейдж
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates bash curl postgresql-client

# Устанавливаем goose
RUN curl -L https://github.com/pressly/goose/releases/latest/download/goose_linux_x86_64 --output /usr/local/bin/goose \
    && chmod +x /usr/local/bin/goose

COPY --from=builder /app/app /app/app
COPY migrations /app/migrations
COPY entrypoint.sh /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
