# Stage 1: сборка приложения и установка утилиты goose
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Устанавливаем необходимые утилиты
RUN apk --no-cache add git

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Устанавливаем goose (он будет собран и сохранён в /go/bin)
RUN go install github.com/pressly/goose/cmd/goose@latest

# Собираем бинарный файл приложения (предполагается, что точка входа находится в cmd)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o victa ./cmd

# Stage 2: финальный образ
FROM alpine:latest
WORKDIR /app

# Устанавливаем сертификаты
RUN apk --no-cache add ca-certificates

# Копируем бинарник приложения и утилиту goose из сборочного образа
COPY --from=builder /app/victa .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копируем миграции (они лежат в internal/database/migrations)
COPY --from=builder /app/internal/migrations /app/migrations

# Копируем entrypoint.sh (скрипт, который запустит миграции, а затем приложение)
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/app/entrypoint.sh"]