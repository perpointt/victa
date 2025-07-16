package webhook

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	"victa/internal/bot/notification_bot"
	"victa/internal/logger"
	"victa/internal/webhook/webhook_common"

	"victa/internal/bot/bot_common"
	"victa/internal/service"
)

type CodemagicWebhookHandler struct {
	*webhook_common.BaseWebhook
	companySvc   *service.CompanyService
	codemagicSvc *service.CodemagicService
}

func NewCodemagicWebhookHandler(
	factory *bot_common.BotFactory,
	logger logger.Logger,
	jwtSvc *service.JWTService,
	companySvc *service.CompanyService,
	codemagicSvc *service.CodemagicService,
) *CodemagicWebhookHandler {
	base := webhook_common.NewBaseWebhook(factory, logger, jwtSvc)
	return &CodemagicWebhookHandler{
		BaseWebhook:  base,
		codemagicSvc: codemagicSvc,
		companySvc:   companySvc,
	}
}

func (h *CodemagicWebhookHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	companyID, err := h.Authorize(c)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var payload struct {
		BuildID string `json:"build_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.SendNewResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	integration, err := h.companySvc.GetCompanyIntegrationByID(ctx, companyID) // +ctx
	if err != nil {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if integration == nil {
		h.SendNewResponse(c, http.StatusBadRequest, "company not found")
		return
	}

	cmCtx, cancelCM := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCM()

	build, err := h.codemagicSvc.GetBuildByID(cmCtx, payload.BuildID, *integration.CodemagicAPIKey) // +ctx
	if err != nil {
		h.SendNewResponse(c, http.StatusBadGateway, err.Error())
		return
	}

	for i, art := range build.Build.Artefacts {
		if strings.EqualFold(art.Type, "apk") {
			url, err := h.codemagicSvc.GetArtifactPublicURL(cmCtx, art.Path, *integration.CodemagicAPIKey)
			if err != nil {
				h.Logger.Warn("codemagic public URL: %v", err)
				break
			}
			build.Build.Artefacts[i].PublicURL = url
			break
		}
	}

	baseBot, err := h.BotFactory.GetBaseBot(*integration.NotificationBotToken, h.Logger)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	bot, err := notification_bot.NewBot(baseBot, *integration.DeployNotificationChatID)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot.SendDeployNotification(build.Application, build.Build)

	h.SendNewResponse(c, http.StatusOK, "OK")
}
