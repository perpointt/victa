package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все настройки приложения (токен бота, параметры БД)
type Config struct {
	// Telegram-бот
	TelegramAdminUserId string // токен бота
	// Telegram-бот
	TelegramToken string // токен бота

	// Параметры подключения к базе
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

// LoadConfig загружает настройки из .env, завершает работу если переменной нет
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла: файл .env обязателен")
	}
	return &Config{
		TelegramAdminUserId: getEnv("TELEGRAM_ADMIN_USER_ID"),
		TelegramToken:       getEnv("TELEGRAM_TOKEN"),
		DBUser:              getEnv("DB_USER"),
		DBPassword:          getEnv("DB_PASSWORD"),
		DBName:              getEnv("DB_NAME"),
		DBHost:              getEnv("DB_HOST"),
		DBPort:              getEnv("DB_PORT"),
	}
}

// GetDbDSN формирует DSN для Postgres
func (c *Config) GetDbDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	log.Fatalf("Переменная окружения %s должна быть установлена", key)
	return ""
}
