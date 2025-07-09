package webhook

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

	integration, err := h.companySvc.GetCompanyIntegrationByID(companyID)
	if err != nil {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if integration == nil {
		h.SendNewResponse(c, http.StatusBadRequest, "company not found")
		return
	}

	buildResp, err := h.codemagicSvc.GetBuildByID(payload.BuildID, *integration.CodemagicAPIKey)
	if err != nil {
		h.SendNewResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	for idx, art := range buildResp.Build.Artefacts {
		if strings.EqualFold(art.Type, "apk") {
			publicURL, err := h.codemagicSvc.GetArtifactPublicURL(
				art.Path, *integration.CodemagicAPIKey,
			)
			if err != nil {
				h.Logger.Errorf(err.Error())
				break
			}
			buildResp.Build.Artefacts[idx].PublicURL = publicURL
			break
		}
	}

	baseBot, err := h.BotFactory.GetBaseBot(*integration.NotificationBotToken, h.Logger)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot, err := notification_bot.NewBot(baseBot, *integration.NotificationChatID)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot.SendNewNotification(buildResp.Application, buildResp.Build)

	h.SendNewResponse(c, http.StatusOK, "OK")
}
