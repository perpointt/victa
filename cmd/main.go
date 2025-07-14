package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"victa/internal/bot/bot_common"
	"victa/internal/bot/victa_bot"
	"victa/internal/config"
	"victa/internal/db"
	"victa/internal/logger"
	"victa/internal/repository"
	"victa/internal/service"
	"victa/internal/webhook"
)

func main() {
	logg := logger.New()
	logg.Info("Приложение запущено")

	cfg := config.LoadConfig(logg)

	conn, err := db.New(cfg.GetDbDSN())
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			logg.Errorf(err.Error())
		}
	}(conn)

	secretBytes := []byte(cfg.InviteSecret)

	var botFactory = bot_common.NewBotFactory()

	userRepo := repository.NewPostgresUserRepo(conn)
	companyRepo := repository.NewPostgresCompanyRepo(conn)
	userCompanyRepo := repository.NewPostgresUserCompanyRepository(conn)
	integrationRepo := repository.NewPostgresCompanyIntegrationRepo(conn)
	appRepo := repository.NewPostgresAppRepo(conn)

	userSvc := service.NewUserService(userRepo, userCompanyRepo)
	companySvc := service.NewCompanyService(companyRepo, integrationRepo)
	inviteSvc := service.NewInviteService(secretBytes)
	appSvc := service.NewAppService(appRepo)
	jwtSvc := service.NewJWTService(cfg.JwtSecret)
	codemagicSvc := service.NewCodemagicService(cfg.CodemagicAPIHost)

	base, err := botFactory.GetBaseBot(cfg.TelegramToken, logg)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота: %v", err)
	}

	b := victa_bot.NewBot(
		base,
		userSvc,
		companySvc,
		inviteSvc,
		appSvc,
		jwtSvc,
	)

	go func() {
		logg.Info("Бот запущен")
		b.Run()
		logg.Warn("Бот остановлен")
	}()

	switch cfg.ENV {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	codemagicHandler := webhook.NewCodemagicWebhookHandler(
		botFactory,
		logg,
		jwtSvc,
		companySvc,
		codemagicSvc,
	)

	gitlabHandler := webhook.NewGitlabWebhookHandler(
		botFactory,
		logg,
		jwtSvc,
		companySvc,
	)

	bugsnagHandler := webhook.NewBugsnagWebhookHandler(
		botFactory,
		logg,
		jwtSvc,
		companySvc,
	)

	r.POST("/webhook/codemagic", codemagicHandler.Handle)
	r.POST("/webhook/gitlab", gitlabHandler.Handle)
	r.POST("/webhook/bugsnag", bugsnagHandler.Handle)

	addr := ":" + cfg.APIPort
	logg.Info("API сервер слушается на %s…", addr)
	if err := r.Run(addr); err != nil {
		logg.Errorf(err.Error())
	}

}
