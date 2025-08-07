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
	releasecron "victa/internal/cron/release"
	"victa/internal/crypto"

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

	utils, err := initUtils()
	if err != nil {
		return err
	}

	services := initServices(cfg, repos, utils)

	//creds, err := os.ReadFile("pc-api-4681420226366194060-658-ad4de3d61ccc.json")
	//if err != nil {
	//	log.Fatalf("read credentials file: %v", err)
	//}

	////// 2. Инициализируем PlayStoreService, передаём содержимое JSON.
	//playSvc, err := service.NewPlayStoreService(ctx, creds)
	//if err != nil {
	//	log.Fatalf("init PlayStoreService: %v", err)
	//}
	//
	//rel, err := playSvc.GetRelease(ctx, "stun.apps.mirror")
	//if err != nil {
	//	log.Fatalf("GetProductionRelease: %v", err)
	//}
	//
	//reviews, err := playSvc.ListReviews(ctx, "stun.apps.ruler", "34545")
	//if err != nil {
	//	log.Fatalf("GetProductionRelease: %v", err)
	//}

	//creds, err := os.ReadFile("AuthKey_9X7C787G6P.p8")
	//if err != nil {
	//	log.Fatalf("read credentials file: %v", err)
	//}
	//
	//appstoreConfig := service.AppStoreConfig{
	//	KeyID:      "9X7C787G6P",
	//	IssuerID:   "f4c08bb6-7921-43f2-96c8-eed8635f48f4",
	//	PrivatePEM: creds,
	//	BaseURL:    cfg.AppStoreAPIHost,
	//}
	//
	//asc, _ := service.NewAppStoreService(appstoreConfig)
	//
	//rel, err := asc.GetRelease(ctx, "6740744282")
	//if err != nil {
	//	log.Fatalf("GetProductionRelease: %v", err)
	//}
	//reviews, err := asc.ListReviews(ctx, "6740744282", "")
	//if err != nil {
	//	log.Fatalf("GetProductionRelease: %v", err)
	//}
	//
	//logg.Info("pkgInfo: %v", rel.Code)
	//logg.Info("pkgInfo: %v", rel.Semantic)
	//logg.Info("pkgInfo: %v", reviews)

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

	botFactory := bot_common.NewBotFactory()

	versionCron := releasecron.NewVersionCron(services.Company, services.App, logg, botFactory, cfg.AppStoreAPIHost)

	versionCron.Run(botCtx)

	g.Go(func() error {
		return versionCron.Start(botCtx)
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

type Utils struct {
	Encryptor *crypto.Encryptor
}

func initUtils() (Utils, error) {
	must := func(v any, err error) (any, error) {
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	encryptor, err := must(crypto.NewEncryptor("master.key"))
	if err != nil {
		return Utils{}, err
	}

	return Utils{
		Encryptor: encryptor.(*crypto.Encryptor),
	}, nil
}

type Repos struct {
	User          *postgres.UserRepo
	Company       *postgres.CompanyRepo
	UserCompany   *postgres.UserCompanyRepo
	CompanySecret *postgres.CompanySecretRepo
	App           *postgres.AppRepo
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
	secret, err := must(postgres.NewCompanySecretRepo(conn))
	if err != nil {
		return Repos{}, err
	}
	app, err := must(postgres.NewAppRepo(conn))
	if err != nil {
		return Repos{}, err
	}

	return Repos{
		User:          user.(*postgres.UserRepo),
		Company:       company.(*postgres.CompanyRepo),
		UserCompany:   userCompany.(*postgres.UserCompanyRepo),
		CompanySecret: secret.(*postgres.CompanySecretRepo),
		App:           app.(*postgres.AppRepo),
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

func initServices(cfg *config.Config, r Repos, u Utils) Services {
	return Services{
		User:      service.NewUserService(r.User, r.UserCompany),
		Company:   service.NewCompanyService(r.Company, r.CompanySecret, *u.Encryptor),
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
