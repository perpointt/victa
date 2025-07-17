package webhook

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"victa/internal/bot/bot_common"
	"victa/internal/bot/notification_bot"
	"victa/internal/domain"
	appErr "victa/internal/errors"
	"victa/internal/logger"
	"victa/internal/service"
	"victa/internal/webhook/webhook_common"
)

type GitlabIssueWebhookHandler struct {
	*webhook_common.BaseWebhook
	companySvc *service.CompanyService
}

func NewGitlabWebhookHandler(
	factory *bot_common.BotFactory,
	logger logger.Logger,
	jwtSvc *service.JWTService,
	companySvc *service.CompanyService,
) *GitlabIssueWebhookHandler {
	base := webhook_common.NewBaseWebhook(factory, logger, jwtSvc)
	return &GitlabIssueWebhookHandler{
		BaseWebhook: base,
		companySvc:  companySvc,
	}
}

func (h *GitlabIssueWebhookHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

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

	if payload.ObjectKind == "issue" &&
		payload.ObjectAttributes.Action == "update" &&
		payload.Changes.ClosedAt != nil &&
		((payload.Changes.ClosedAt.Previous == nil && payload.Changes.ClosedAt.Current != nil) ||
			(payload.Changes.ClosedAt.Current == nil && payload.Changes.ClosedAt.Previous != nil)) {

		h.SendNewResponse(c, http.StatusOK, "OK, but ignored")
		return
	}

	botToken, err := h.companySvc.GetCompanySecret(ctx, companyID, domain.SecretNotificationBotToken)
	if err != nil && !errors.Is(err, appErr.ErrSecretNotFound) {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if botToken == nil {
		h.SendNewResponse(c, http.StatusBadRequest, "notification bot token not found")
		return
	}

	chatID, err := h.companySvc.GetCompanySecret(ctx, companyID, domain.SecretIssuesNotificationChatID)
	if err != nil && !errors.Is(err, appErr.ErrSecretNotFound) {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if chatID == nil {
		h.SendNewResponse(c, http.StatusBadRequest, "notification chat id not found")
		return
	}

	baseBot, err := h.BotFactory.GetBaseBot(string(botToken), h.Logger)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	bot, err := notification_bot.NewBot(baseBot, string(chatID))
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot.SendIssueNotification(payload)

	h.SendNewResponse(c, http.StatusOK, "OK")
}
