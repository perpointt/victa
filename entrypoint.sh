#!/bin/sh
set -e

echo "Запуск миграций базы данных..."

# Формируем строку подключения для goose
# Пример DSN: host=localhost port=5433 user=victa dbname=victa_db password=FDTpmXuaKVxB584Q9y6fGP sslmode=disable
DB_DSN="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} dbname=${DB_NAME} password=${DB_PASSWORD} sslmode=disable"

if [ -z "$DB_DSN" ]; then
  echo "Переменная DB_DSN не установлена"
  exit 1
fi

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="$DB_DSN"

# Запускаем миграции. Файлы миграций должны находиться в /app/migrations.
goose -dir /app/migrations up

echo "Миграции завершены, запускаем приложение..."
exec /app/victa