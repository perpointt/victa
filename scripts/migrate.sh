#!/bin/sh
set -e

# Загружаем переменные из файла .env
if [ -f .env ]; then
  # Экспортируем все переменные, исключая комментарии
  export $(grep -v '^#' .env | xargs)
else
  echo ".env file not found!"
  exit 1
fi

# Формируем строку подключения для goose
# Пример DSN: host=localhost port=5433 user=victa dbname=victa_db password=FDTpmXuaKVxB584Q9y6fGP sslmode=disable
DB_DSN="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable"

# Указываем директорию миграций
MIGRATIONS_DIR="./migrations"

echo "Running goose migrations with DSN: ${DB_DSN}"
goose -dir "${MIGRATIONS_DIR}" postgres "${DB_DSN}" up
