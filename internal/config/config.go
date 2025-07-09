package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"victa/internal/logger"
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

func LoadConfig(logger logger.Logger) *Config {
	if err := godotenv.Load(); err != nil {
		logger.Errorf(err.Error())
	}
	return &Config{
		TelegramToken:    getEnv("TELEGRAM_TOKEN", logger),
		TelegramBotName:  getEnv("TELEGRAM_BOT_NAME", logger),
		InviteSecret:     getEnv("INVITE_SECRET", logger),
		JwtSecret:        getEnv("JWT_SECRET", logger),
		DBUser:           getEnv("DB_USER", logger),
		DBPassword:       getEnv("DB_PASSWORD", logger),
		DBName:           getEnv("DB_NAME", logger),
		DBHost:           getEnv("DB_HOST", logger),
		DBPort:           getEnv("DB_PORT", logger),
		APIPort:          getEnv("API_PORT", logger),
		CodemagicAPIHost: getEnv("CODEMAGIC_API_HOST", logger),
		ENV:              getEnv("ENV", logger),
	}
}

func (c *Config) GetDbDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key string, logger logger.Logger) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	logger.Warn("Переменная окружения %s должна быть установлена", key)
	return ""
}
