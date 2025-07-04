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

	secretBytes := []byte(cfg.InviteSecret)

	userRepo := repository.NewPostgresUserRepo(conn)
	companyRepo := repository.NewPostgresCompanyRepo(conn)

	userSvc := service.NewUserService(userRepo)
	companySvc := service.NewCompanyService(companyRepo)
	inviteSvc := service.NewInviteService(secretBytes)

	// Инициализируем Telegram-бота
	b, err := bot.NewBot(*cfg, userSvc, companySvc, inviteSvc)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота: %v", err)
	}

	log.Println("Bot started...")
	b.Run()
}
