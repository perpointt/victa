package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"victa/internal/bot/bot_common"
	"victa/internal/bot/victa_bot"
	"victa/internal/config"
	"victa/internal/db"
	"victa/internal/logger"
	"victa/internal/repository/postgres"
	"victa/internal/service"
	"victa/internal/webhook"
)

func main() {
	logg := logger.New(nil, logger.LevelInfo) // stderr, INFO и выше

	if err := run(logg); err != nil {
		logg.Error("shutdown: %v", err)
		os.Exit(1)
	}
}

func run(logg logger.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn, err := initDB(ctx, cfg.GetDbDSN())
	if err != nil {
		return err
	}
	defer func() {
		_ = dbConn.Close()
	}()

	repos, err := initRepos(dbConn)
	if err != nil {
		return err
	}
	services := initServices(cfg, repos)

	botBase, err := bot_common.NewBotFactory().GetBaseBot(cfg.TelegramToken, logg)
	if err != nil {
		return fmt.Errorf("init telegram bot: %w", err)
	}
	tgBot := victa_bot.New(
		botBase,
		cfg.TelegramBotName,
		services.User,
		services.Company,
		services.Invite,
		services.App,
		services.JWT,
	)

	router := buildRouter(cfg, logg, services)

	srv := &http.Server{
		Addr:         ":" + cfg.APIPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logg.Info("HTTP‑API слушает %s", srv.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	botCtx, cancelBot := context.WithCancel(ctx)

	g.Go(func() error {
		logg.Info("Telegram‑бот запущен")
		return tgBot.Run(botCtx)
	})

	g.Go(func() error {
		<-gCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		cancelBot()
		return nil
	})

	return g.Wait()
}

func initDB(ctx context.Context, dsn string) (*sql.DB, error) {
	conn, err := db.New(dsn)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}
	return conn, nil
}

type Repos struct {
	User        *postgres.UserRepo
	Company     *postgres.CompanyRepo
	UserCompany *postgres.UserCompanyRepo
	Integration *postgres.CompanyIntegrationRepo
	App         *postgres.AppRepo
}

func initRepos(conn *sql.DB) (Repos, error) {
	must := func(v any, err error) (any, error) {
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	user, err := must(postgres.NewUserRepo(conn))
	if err != nil {
		return Repos{}, err
	}
	company, err := must(postgres.NewCompanyRepo(conn))
	if err != nil {
		return Repos{}, err
	}
	userCompany, err := must(postgres.NewUserCompanyRepo(conn))
	if err != nil {
		return Repos{}, err
	}
	integration, err := must(postgres.NewCompanyIntegrationRepo(conn))
	if err != nil {
		return Repos{}, err
	}
	app, err := must(postgres.NewAppRepo(conn))
	if err != nil {
		return Repos{}, err
	}

	return Repos{
		User:        user.(*postgres.UserRepo),
		Company:     company.(*postgres.CompanyRepo),
		UserCompany: userCompany.(*postgres.UserCompanyRepo),
		Integration: integration.(*postgres.CompanyIntegrationRepo),
		App:         app.(*postgres.AppRepo),
	}, nil
}

type Services struct {
	User      *service.UserService
	Company   *service.CompanyService
	Invite    *service.InviteService
	App       *service.AppService
	JWT       *service.JWTService
	Codemagic *service.CodemagicService
}

func initServices(cfg *config.Config, r Repos) Services {
	return Services{
		User:      service.NewUserService(r.User, r.UserCompany),
		Company:   service.NewCompanyService(r.Company, r.Integration),
		Invite:    service.NewInviteService([]byte(cfg.InviteSecret), 48*time.Hour),
		App:       service.NewAppService(r.App),
		JWT:       service.NewJWTService(cfg.JwtSecret),
		Codemagic: service.NewCodemagicService(cfg.CodemagicAPIHost),
	}
}

func buildRouter(cfg *config.Config, logg logger.Logger, s Services) *gin.Engine {
	if cfg.ENV == "prod" || cfg.ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	botFactory := bot_common.NewBotFactory()

	r.POST("/webhook/codemagic",
		webhook.NewCodemagicWebhookHandler(botFactory, logg, s.JWT, s.Company, s.Codemagic).Handle,
	)
	r.POST("/webhook/gitlab",
		webhook.NewGitlabWebhookHandler(botFactory, logg, s.JWT, s.Company).Handle,
	)
	r.POST("/webhook/bugsnag",
		webhook.NewBugsnagWebhookHandler(botFactory, logg, s.JWT, s.Company).Handle,
	)

	return r
}
