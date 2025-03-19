package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все необходимые настройки приложения victa.
type Config struct {
	Port          string // Порт, на котором будет запущен сервер
	JWTSecret     string // Секрет для подписи JWT токенов
	EncryptionKey string // Ключ для шифрования чувствительных данных (32 байта для AES-256)
	Env           string // Окружение (например, dev или prod)

	// Настройки подключения к базе данных
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

// LoadConfig загружает конфигурацию из .env файла. Если .env отсутствует или переменная не задана, приложение завершится с ошибкой.
func LoadConfig() *Config {
	// Загружаем переменные из .env файла.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла: файл .env обязателен")
	}

	return &Config{
		Port:          getEnv("PORT"),
		JWTSecret:     getEnv("JWT_SECRET"),
		EncryptionKey: getEnv("ENCRYPTION_KEY"),
		Env:           getEnv("ENV"),
		DBUser:        getEnv("DB_USER"),
		DBPassword:    getEnv("DB_PASSWORD"),
		DBName:        getEnv("DB_NAME"),
		DBHost:        getEnv("DB_HOST"),
		DBPort:        getEnv("DB_PORT"),
	}
}

// GetDbDSN формирует строку подключения к базе данных.
func (c *Config) GetDbDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// getEnv возвращает значение переменной окружения или завершает работу, если переменная не установлена.
func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Переменная окружения %s должна быть установлена", key)
	return ""
}
