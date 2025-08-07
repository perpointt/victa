package releasecron

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"victa/internal/domain"
	"victa/internal/logger"
	"victa/internal/service"
)

const (
	scheduleHourly = "@hourly"
	companyPool    = 10
	appPool        = 5
	requestTimeout = 30 * time.Second
)

var (
	playIDRe = regexp.MustCompile(`\bid=([^&]+)`)
	appIDRe  = regexp.MustCompile(`/id(\d+)\b`)
)

// CronVersion шедулит проверку версий.
type CronVersion struct {
	companySvc      *service.CompanyService
	appSvc          *service.AppService
	logger          logger.Logger
	appStoreBaseURL string
}

// NewVersionCron создаёт cron, принимая базовый URL App Store API.
func NewVersionCron(
	companySvc *service.CompanyService,
	appSvc *service.AppService,
	logger logger.Logger,
	appStoreBaseURL string,
) *CronVersion {
	return &CronVersion{
		companySvc:      companySvc,
		appSvc:          appSvc,
		logger:          logger,
		appStoreBaseURL: appStoreBaseURL,
	}
}

// Start запускает cron и блокирует до ctx.Done().
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
	// Google Play
	googleKey, _ := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretGoogleJSON)

	// App Store
	appleP8, errP8 := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretAppleP8)
	issuerID, errIss := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretAppleIssuerID)
	keyID, errKey := c.companySvc.GetSecret(ctx, comp.ID, domain.SecretAppleKeyID)
	if errP8 != nil || errIss != nil || errKey != nil {
		appleP8 = nil
	}

	apps, err := c.appSvc.GetAllByCompanyID(ctx, comp.ID)
	if err != nil {
		c.logger.Error("version-cron[%s]: GetAllByCompanyID: %v", comp.Name, err)
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
				ctx, comp.Name, app,
				googleKey,
				appleP8, string(issuerID), string(keyID),
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
) error {
	ctxReq, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	// Google Play
	if googleKey != nil && app.PlayStoreURL != nil {
		if pkg := parsePlayID(*app.PlayStoreURL); pkg != "" {
			playServiceSvc, err := service.NewPlayStoreService(ctxReq, googleKey)
			if err != nil {
				c.logger.Error("PlayService init: %v", err)
			} else if rel, err := playServiceSvc.GetRelease(ctxReq, pkg); err != nil {
				c.logger.Error("Play[%s/%s]: %v", companyName, pkg, err)
			} else {
				fmt.Printf(
					"Company=%s, App=%s (Play) → %s (%d)\n",
					companyName, app.Name,
					rel.Semantic, rel.Code,
				)
			}
		}
	}

	// App Store
	if appleP8 != nil && app.AppStoreURL != nil {
		if appID := parseAppStoreID(*app.AppStoreURL); appID != "" {
			cfg := service.AppStoreConfig{
				KeyID:      appleKeyID,
				IssuerID:   appleIssuerID,
				PrivatePEM: appleP8,
				BaseURL:    c.appStoreBaseURL, // используем переданный в конструктор
			}
			appstoreSvc, err := service.NewAppStoreService(cfg)
			if err != nil {
				c.logger.Error("AppStoreService init: %v", err)
			} else if rel, err := appstoreSvc.GetRelease(ctxReq, appID); err != nil {
				c.logger.Error("AppStore[%s/%s]: %v", companyName, appID, err)
			} else {
				fmt.Printf(
					"Company=%s, App=%s (AppStore) → %s (%d)\n",
					companyName, app.Name,
					rel.Semantic, rel.Code,
				)
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
