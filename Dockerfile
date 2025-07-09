# ————— Stage 1: Build —————
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Утилиты
RUN apk add --no-cache git

# Модули
COPY go.mod go.sum ./
RUN go mod download

# Код и миграции
COPY . .

# Установим goose для миграций
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o victa ./cmd

# ————— Stage 2: Runtime —————
FROM alpine:latest
WORKDIR /app

# Системные зависимости
RUN apk add --no-cache ca-certificates curl

# Копируем артефакты из builder
COPY --from=builder /app/victa .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/migrations ./migrations
COPY --from=builder /app/entrypoint.sh .

# Скрипт-запускатель
RUN chmod +x entrypoint.sh

# Порт приложения
EXPOSE 3000

ENTRYPOINT ["./entrypoint.sh"]