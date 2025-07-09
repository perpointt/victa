#!/bin/bash
set -e

# Определяем директории для релиза, бэкапа и текущего проекта.
RELEASE_DIR="/opt/victa/releases/$(date +%Y%m%d%H%M%S)"
BACKUP_DIR="/opt/victa/backup"
PROJECT_DIR="/opt/victa/current"

echo "Создаем директорию релиза: $RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

echo "Создаем бэкап текущей версии..."
if [ -d "$BACKUP_DIR" ]; then
  echo "Удаляем предыдущий бэкап: $BACKUP_DIR"
  rm -rf "$BACKUP_DIR"
fi

if [ -d "$PROJECT_DIR" ]; then
  mv "$PROJECT_DIR" "$BACKUP_DIR"
  echo "Бэкап выполнен, текущая версия перемещена в: $BACKUP_DIR"
fi

echo "Распаковываем новый архив..."
tar xzf /opt/victa/releases/victa.tar.gz -C "$RELEASE_DIR"

echo "Обновляем симлинк: $PROJECT_DIR -> $RELEASE_DIR"
ln -sfn "$RELEASE_DIR" "$PROJECT_DIR"

echo "Перезапускаем контейнеры..."
cd "$PROJECT_DIR"

docker-compose -p victa -f docker-compose.yaml down
docker-compose -p victa -f docker-compose.yaml up -d --build

echo "Деплой завершен."