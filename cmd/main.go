package main

import (
	"log"
	"net/http"
	"victa/internal/bot"
	"victa/internal/config"
	"victa/internal/db"
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
	cmSvc := service.NewCodemagicService(cfg.CodemagicAPIHost)

	b, err := bot.NewBot(
		*cfg,
		userSvc,
		companySvc,
		inviteSvc,
		appSvc,
		jwtSvc,
	)
	if err != nil {
		log.Fatalf("Ошибка при инициализации бота: %v", err)
	}

	go func() {
		log.Println("Bot started…")
		b.Run()
		log.Println("Bot stopped")
	}()

	buildHandler := webhook.NewBuildWebhookHandler(jwtSvc, companySvc, cmSvc)
	http.Handle("/webhook/build", buildHandler)

	addr := ":" + cfg.APIPort
	log.Printf("API server listening on %s…", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
