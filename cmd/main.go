package main

import (
	"log"
	"victa/internal/bot"
	"victa/internal/config"
	"victa/internal/db"
	"victa/internal/repository"
	"victa/internal/service"
)

func main() {
	// Загружаем конфиг
	cfg := config.LoadConfig()

	// Подключаемся к БД
	conn, err := db.New(cfg.GetDbDSN())
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer conn.Close()

	userRepo := repository.NewPostgresUserRepo(conn)
	userSettingsRepo := repository.NewPostgresUserSettingsRepo(conn)
	companyRepo := repository.NewPostgresCompanyRepo(conn)

	userSvc := service.NewUserService(userRepo)
	userSettingsSvc := service.NewUserSettingsService(userSettingsRepo)
	companySvc := service.NewCompanyService(companyRepo)

	// Инициализируем Telegram-бота
	b, err := bot.NewBot(*cfg, userSvc, userSettingsSvc, companySvc)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота: %v", err)
	}

	log.Println("Bot started...")
	b.Run()
}
