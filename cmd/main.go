package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"victa/internal/config"
	"victa/internal/db"
	"victa/internal/new_bot/bot_common"
	"victa/internal/new_bot/victa_bot"
	"victa/internal/repository"
	"victa/internal/service"
	"victa/internal/webhook"
)

func main() {
	cfg := config.LoadConfig()

	conn, err := db.New(cfg.GetDbDSN())
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer conn.Close()

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

	base, err := botFactory.GetBaseBot(cfg.TelegramToken)
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
		log.Println("Bot started…")
		b.Run()
		log.Println("Bot stopped")
	}()

	switch cfg.ENV {
	case "prod", "production":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	buildHandler := webhook.NewBuildWebhookHandler(botFactory, jwtSvc, companySvc, codemagicSvc)
	r.POST("/webhook/build", buildHandler.Handle)

	addr := ":" + cfg.APIPort
	log.Printf("API server listening on %s…", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}

}
