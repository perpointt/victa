package webhook

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"victa/internal/bot/bot_common"
	"victa/internal/bot/notification_bot"
	"victa/internal/domain"
	"victa/internal/logger"
	"victa/internal/service"
	"victa/internal/webhook/webhook_common"
)

type GitlabWebhookHandler struct {
	*webhook_common.BaseWebhook
	companySvc *service.CompanyService
}

func NewGitlabWebhookHandler(
	factory *bot_common.BotFactory,
	logger logger.Logger,
	jwtSvc *service.JWTService,
	companySvc *service.CompanyService,
) *GitlabWebhookHandler {
	base := webhook_common.NewBaseWebhook(factory, logger, jwtSvc)
	return &GitlabWebhookHandler{
		BaseWebhook: base,
		companySvc:  companySvc,
	}
}

func (h *GitlabWebhookHandler) Handle(c *gin.Context) {
	companyID, err := h.Authorize(c)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var payload domain.GitlabWebhook

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

	baseBot, err := h.BotFactory.GetBaseBot(*integration.NotificationBotToken, h.Logger)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot, err := notification_bot.NewBot(baseBot, *integration.IssuesNotificationChatID)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot.SendIssueNotification(payload)

	h.SendNewResponse(c, http.StatusOK, "OK")
}
