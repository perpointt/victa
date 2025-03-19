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
	cfg := config.LoadConfig()

	// Устанавливаем режим Gin в зависимости от ENV
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	db, err := sql.Open("postgres", cfg.GetDbDSN())
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Инициализация для компаний
	companyRepo := repository.NewCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepo)
	// Передаем companyService и userService в NewCompanyHandler
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	companyHandler := handler.NewCompanyHandler(companyService, userService)

	// Инициализация для пользователей
	userHandler := handler.NewUserHandler(userService)

	// Инициализация для приложений
	appRepo := repository.NewAppRepository(db)
	appService := service.NewAppService(appRepo)
	appHandler := handler.NewAppHandler(appService)

	// Инициализация для связи пользователей и компаний
	userCompanyRepo := repository.NewUserCompanyRepository(db)

	// Инициализация для аутентификации
	authService := service.NewAuthService(userRepo, companyRepo, userCompanyRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Настройка маршрутов с подключенной JWT-миддлварой
	r := router.SetupRouter(companyHandler, userHandler, appHandler, authHandler, cfg.JWTSecret)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
