package webhook

import (
	"errors"
	"net/http"
	"victa/internal/domain"
	appErr "victa/internal/errors"

	"github.com/gin-gonic/gin"
	"victa/internal/bot/bot_common"
	"victa/internal/bot/notification_bot"
	"victa/internal/logger"
	"victa/internal/service"
	"victa/internal/webhook/webhook_common"
)

type BugsnagWebhookHandler struct {
	*webhook_common.BaseWebhook
	companySvc *service.CompanyService
}

func NewBugsnagWebhookHandler(
	factory *bot_common.BotFactory,
	logger logger.Logger,
	jwtSvc *service.JWTService,
	companySvc *service.CompanyService,
) *BugsnagWebhookHandler {
	base := webhook_common.NewBaseWebhook(factory, logger, jwtSvc)
	return &BugsnagWebhookHandler{
		BaseWebhook: base,
		companySvc:  companySvc,
	}
}
func (h *BugsnagWebhookHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()

	companyID, err := h.Authorize(c)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	var payload domain.BugsnagWebhook
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.SendNewResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	botToken, err := h.companySvc.GetSecret(ctx, companyID, domain.SecretNotificationBotToken)
	if err != nil && !errors.Is(err, appErr.ErrSecretNotFound) {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if botToken == nil {
		h.SendNewResponse(c, http.StatusBadRequest, "notification bot token not found")
		return
	}

	chatID, err := h.companySvc.GetSecret(ctx, companyID, domain.SecretErrorsNotificationChatID)
	if err != nil && !errors.Is(err, appErr.ErrSecretNotFound) {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if chatID == nil {
		h.SendNewResponse(c, http.StatusBadRequest, "notification chat id not found")
		return
	}

	baseBot, err := h.BotFactory.GetBaseBot(
		string(botToken),
		h.Logger,
	)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot, err := notification_bot.NewBot(baseBot, string(chatID))
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot.SendBugsnagNotification(payload)

	h.SendNewResponse(c, http.StatusOK, "OK")
}
