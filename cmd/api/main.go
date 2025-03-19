package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"victa/internal/config"
	"victa/internal/handler"
	"victa/internal/repository"
	"victa/internal/router"
	"victa/internal/service"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Устанавливаем режим Gin в зависимости от переменной ENV
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", cfg.GetDbDSN())
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Инициализируем репозиторий, сервис и обработчик для компаний
	companyRepo := repository.NewCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepo)
	companyHandler := handler.NewCompanyHandler(companyService)

	// Настраиваем маршрутизацию
	r := router.SetupRouter(companyHandler)

	// Запускаем сервер
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
