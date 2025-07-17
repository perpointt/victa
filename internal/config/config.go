package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
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
	AppStoreAPIHost  string
	ENV              string
}

// Load читает переменные окружения и возвращает Config.
// Если .env не найден или отсутствует обязательная переменная — ошибка поднимается наверх.
func Load() (*Config, error) {
	// .env может отсутствовать в продакшене — это не критично
	_ = godotenv.Load()

	var missing []string
	mustEnv := func(k string) string {
		if v, ok := os.LookupEnv(k); ok {
			return v
		}
		missing = append(missing, k)
		return ""
	}

	cfg := &Config{
		TelegramToken:    mustEnv("TELEGRAM_TOKEN"),
		TelegramBotName:  mustEnv("TELEGRAM_BOT_NAME"),
		InviteSecret:     mustEnv("INVITE_SECRET"),
		JwtSecret:        mustEnv("JWT_SECRET"),
		DBUser:           mustEnv("DB_USER"),
		DBPassword:       mustEnv("DB_PASSWORD"),
		DBName:           mustEnv("DB_NAME"),
		DBHost:           mustEnv("DB_HOST"),
		DBPort:           mustEnv("DB_PORT"),
		APIPort:          mustEnv("API_PORT"),
		CodemagicAPIHost: mustEnv("CODEMAGIC_API_HOST"),
		AppStoreAPIHost:  mustEnv("APP_STORE_API_HOST"),
		ENV:              mustEnv("ENV"),
	}

	if len(missing) > 0 {
		return nil, errors.New("missing env vars: " + strings.Join(missing, ", "))
	}
	return cfg, nil
}

// GetDbDSN формирует строку подключения к Postgres.
func (c *Config) GetDbDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}
