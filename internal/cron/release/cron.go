package releasecron

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"victa/internal/bot/bot_common"
	"victa/internal/bot/notification_bot"
	"victa/internal/domain"
	appErr "victa/internal/errors"
	"victa/internal/logger"
	"victa/internal/repository"
	"victa/internal/service"
)

const (
	scheduleHourly = "@hourly"
	companyPool    = 10
	appPool        = 5
	requestTimeout = 30 * time.Second

	// захардкодим ID маркетов (как в таблице app_markets)
	marketGooglePlayID int64 = 1
	marketAppStoreID   int64 = 2
)

var (
	playIDRe = regexp.MustCompile(`\bid=([^&]+)`)
	appIDRe  = regexp.MustCompile(`/id(\d+)\b`)
)

// CronVersion шедулит проверку версий и уведомления.
type CronVersion struct {
	companySvc      *service.CompanyService
	appSvc          *service.AppService
	stateRepo       repository.AppMarketStateRepository
	logger          logger.Logger
	botFactory      *bot_common.BotFactory
	appStoreBaseURL string
}

func NewVersionCron(
	companySvc *service.CompanyService,
	appSvc *service.AppService,
	stateRepo repository.AppMarketStateRepository,
	logger logger.Logger,
	botFactory *bot_common.BotFactory,
	appStoreBaseURL string,
) *CronVersion {
	return &CronVersion{
		companySvc:      companySvc,
		appSvc:          appSvc,
		stateRepo:       stateRepo,
		logger:          logger,
		botFactory:      botFactory,
		appStoreBaseURL: appStoreBaseURL,
	}
}

// Start запускает cron.
func (c *CronVersion) Start(ctx context.Context) error {
	cr := cron.New()
	_, err := cr.AddFunc(scheduleHourly, func() { c.Run(ctx) })
	if err != nil {
		return err
	}
	c.logger.Info("version-cron: старт, расписание %s", scheduleHourly)
	cr.Start()
	<-ctx.Done()
	st := cr.Stop()
	<-st.Done()
	return nil
}

// Run — один проход по всем компаниям и приложениям.
func (c *CronVersion) Run(ctx context.Context) {
	c.logger.Info("version-cron: раунд start")

	comps, err := c.companySvc.GetAll(ctx)
	if err != nil {
		c.logger.Error("version-cron: GetAll companies: %v", err)
		return
	}

	semComp := semaphore.NewWeighted(companyPool)
	eg, ctx := errgroup.WithContext(ctx)

	for _, comp := range comps {
		comp := comp
		eg.Go(func() error {
			if err := semComp.Acquire(ctx, 1); err != nil {
				return err
			}
			defer semComp.Release(1)
			return c.handleCompany(ctx, comp)
		})
	}
	_ = eg.Wait()
	c.logger.Info("version-cron: раунд done")
}

func (c *CronVersion) handleCompany(ctx context.Context, comp domain.Company) error {
	// токен и чат для уведомлений
	tokenB, err := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretNotificationBotToken)
	if err != nil && !errors.Is(err, appErr.ErrSecretNotFound) {
		c.logger.Error("company[%s]: failed to load bot token: %v", comp.Name, err)
		return nil
	}
	if tokenB == nil {
		c.logger.Info("company[%s]: no notification bot token, пропускаем", comp.Name)
		return nil
	}

	chatB, err := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretVersionsNotificationChatID)
	if err != nil && !errors.Is(err, appErr.ErrSecretNotFound) {
		c.logger.Error("company[%s]: failed to load versions chat ID: %v", comp.Name, err)
		return nil
	}
	if chatB == nil {
		c.logger.Info("company[%s]: no versions chat ID, пропускаем", comp.Name)
		return nil
	}

	baseBot, err := c.botFactory.GetBaseBot(string(tokenB), c.logger)
	if err != nil {
		c.logger.Error("company[%s]: botFactory error: %v", comp.Name, err)
		return nil
	}
	notifyBot, err := notification_bot.NewBot(baseBot, string(chatB))
	if err != nil {
		c.logger.Error("company[%s]: NewBot error: %v", comp.Name, err)
		return nil
	}

	// секреты сторов
	googleKey, _ := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretGoogleJSON)

	appleP8, errP8 := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretAppleP8)
	issuerID, errIss := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretAppleIssuerID)
	keyID, errKey := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretAppleKeyID)
	if errP8 != nil || errIss != nil || errKey != nil {
		appleP8 = nil
	}

	apps, err := c.appSvc.GetAllByCompanyID(ctx, comp.ID)
	if err != nil {
		c.logger.Error("company[%s]: GetAllByCompanyID: %v", comp.Name, err)
		return nil
	}

	semApp := semaphore.NewWeighted(appPool)
	eg, _ := errgroup.WithContext(ctx)

	for _, app := range apps {
		app := app
		eg.Go(func() error {
			if err := semApp.Acquire(ctx, 1); err != nil {
				return err
			}
			defer semApp.Release(1)
			return c.handleApp(
				ctx,
				comp.Name,
				app,
				googleKey,
				appleP8, string(issuerID), string(keyID),
				notifyBot,
			)
		})
	}

	_ = eg.Wait()
	return nil
}

func (c *CronVersion) handleApp(
	ctx context.Context,
	companyName string,
	app domain.App,
	googleKey []byte,
	appleP8 []byte,
	appleIssuerID, appleKeyID string,
	notifyBot *notification_bot.Bot,
) error {
	ctxReq, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	// -------- Google Play --------
	if googleKey != nil && app.PlayStoreURL != nil {
		if pkg := parsePlayID(*app.PlayStoreURL); pkg != "" {
			playSvc, err := service.NewPlayStoreService(ctxReq, googleKey, c.stateRepo, 1)
			if err != nil {
				c.logger.Error("Play init: %v", err)
			} else if rel, err := playSvc.CheckNewRelease(
				ctxReq,
				app,
				pkg,
				c.stateRepo,
				marketGooglePlayID,
			); err != nil {
				c.logger.Error("Play[%s/%s]: %v", companyName, pkg, err)
			} else if rel != nil {
				// есть новая версия — шлём уведомление
				notifyBot.SendVersionNotification(*rel)
			}
		}
	}

	// -------- App Store --------
	if appleP8 != nil && app.AppStoreURL != nil {
		if appID := parseAppStoreID(*app.AppStoreURL); appID != "" {
			cfg := service.AppStoreConfig{
				KeyID:      appleKeyID,
				IssuerID:   appleIssuerID,
				PrivatePEM: appleP8,
				BaseURL:    c.appStoreBaseURL,
			}
			appStoreSvc, err := service.NewAppStoreService(cfg, c.stateRepo, 2)
			if err != nil {
				c.logger.Error("AppStore init: %v", err)
			} else if rel, err := appStoreSvc.CheckNewRelease(
				ctxReq,
				app,
				appID,
				"", // bundleID, если захочешь парсить из URL и передавать
				c.stateRepo,
				marketAppStoreID,
			); err != nil {
				c.logger.Error("AppStore[%s/%s]: %v", companyName, appID, err)
			} else if rel != nil {
				notifyBot.SendVersionNotification(*rel)
			}
		}
	}

	return nil
}

func parsePlayID(u string) string {
	if m := playIDRe.FindStringSubmatch(u); len(m) == 2 {
		return m[1]
	}
	return ""
}

func parseAppStoreID(u string) string {
	if m := appIDRe.FindStringSubmatch(u); len(m) == 2 {
		return m[1]
	}
	return ""
}
