#!/bin/bash
set -e

# Пути к текущему проекту, бэкапу и архиву с новой версией
PROJECT_DIR="/opt/victa/current"
BACKUP_DIR="/opt/victa/backup"
TAR_FILE="/opt/victa/victa.tar.gz"

echo "Создаем бэкап текущей версии..."
if [ -d "$BACKUP_DIR" ]; then
  echo "Удаляем старый бэкап: $BACKUP_DIR"
  rm -rf "$BACKUP_DIR"
fi

if [ -d "$PROJECT_DIR" ]; then
  mv "$PROJECT_DIR" "$BACKUP_DIR"
  echo "Бэкап сохранен в: $BACKUP_DIR"
fi

echo "Готовим директорию для новой версии: $PROJECT_DIR"
mkdir -p "$PROJECT_DIR"

echo "Распаковываем новый архив ($TAR_FILE) в $PROJECT_DIR..."
tar xzf "$TAR_FILE" -C "$PROJECT_DIR"

echo "Перезапускаем контейнеры..."
cd "$PROJECT_DIR"

docker builder prune --force
docker-compose -p victa -f docker-compose.yaml down
docker-compose -p victa -f docker-compose.yaml up -d --build

echo "Деплой завершен."
