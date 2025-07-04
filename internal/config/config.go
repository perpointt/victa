package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramAdminUserId string
	TelegramToken       string
	TelegramBotName     string
	InviteSecret        string
	DBUser              string
	DBPassword          string
	DBName              string
	DBHost              string
	DBPort              string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла: файл .env обязателен")
	}
	return &Config{
		TelegramAdminUserId: getEnv("TELEGRAM_ADMIN_USER_ID"),
		TelegramToken:       getEnv("TELEGRAM_TOKEN"),
		TelegramBotName:     getEnv("TELEGRAM_BOT_NAME"),
		InviteSecret:        getEnv("INVITE_SECRET"),
		DBUser:              getEnv("DB_USER"),
		DBPassword:          getEnv("DB_PASSWORD"),
		DBName:              getEnv("DB_NAME"),
		DBHost:              getEnv("DB_HOST"),
		DBPort:              getEnv("DB_PORT"),
	}
}

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
