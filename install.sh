#!/bin/bash
set -e

# Обновляем систему и устанавливаем зависимости, если они отсутствуют.
echo "Обновляем систему..."
sudo apt update

# Проверяем, установлен ли Docker.
if ! command -v docker &> /dev/null; then
    echo "Docker не установлен. Устанавливаем Docker..."
    sudo apt install -y ca-certificates curl gnupg lsb-release
    sudo mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    echo "Docker установлен."
else
    echo "Docker уже установлен."
fi

# Проверяем, установлен ли docker-compose (если используется отдельный бинарный файл)
if ! command -v docker-compose &> /dev/null; then
    echo "docker-compose не найден. Устанавливаем docker-compose..."
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo "docker-compose установлен."
else
    echo "docker-compose уже установлен."
fi

# Строим образы и поднимаем контейнеры в режиме detached (продакшн).
echo "Строим образы и поднимаем контейнеры..."
sudo docker-compose -f docker-compose.prod.yaml up -d --build

echo "Развертывание завершено. Приложение victa запущено."