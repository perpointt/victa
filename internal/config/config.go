package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken    string
	TelegramBotName  string
	InviteSecret     string
	JwtSecret        string
	DBUser           string
	DBPassword       string
	DBName           string
	DBHost           string
	DBPort           string
	APIPort          string
	CodemagicAPIHost string
	ENV              string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла: файл .env обязателен")
	}
	return &Config{
		TelegramToken:    getEnv("TELEGRAM_TOKEN"),
		TelegramBotName:  getEnv("TELEGRAM_BOT_NAME"),
		InviteSecret:     getEnv("INVITE_SECRET"),
		JwtSecret:        getEnv("JWT_SECRET"),
		DBUser:           getEnv("DB_USER"),
		DBPassword:       getEnv("DB_PASSWORD"),
		DBName:           getEnv("DB_NAME"),
		DBHost:           getEnv("DB_HOST"),
		DBPort:           getEnv("DB_PORT"),
		APIPort:          getEnv("API_PORT"),
		CodemagicAPIHost: getEnv("CODEMAGIC_API_HOST"),
		ENV:              getEnv("ENV"),
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
