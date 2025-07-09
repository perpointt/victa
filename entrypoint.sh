#!/bin/sh
set -e

echo "Запуск миграций базы данных..."

if [ -z "$DATABASE_DSN" ]; then
  echo "Переменная DATABASE_DSN не установлена"
  exit 1
fi

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="$DATABASE_DSN"

# Запускаем миграции. Файлы миграций должны находиться в /app/migrations.
goose -dir /app/migrations up

echo "Миграции завершены, запускаем приложение..."
exec /app/victa