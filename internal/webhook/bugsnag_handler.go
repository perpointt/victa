package webhook

import (
	"bytes"
	"io"
	"net/http"

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
	companyID, err := h.Authorize(c)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Читаем и логируем тело запроса
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.Logger.Errorf("Failed to read request body: %v", err)
		h.SendNewResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	h.Logger.Info("Incoming Bugsnag payload: %s", string(bodyBytes))

	// Получаем интеграцию для компании
	integration, err := h.companySvc.GetCompanyIntegrationByID(companyID)
	if err != nil {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if integration == nil {
		h.SendNewResponse(c, http.StatusBadRequest, "company not found")
		return
	}

	// Создаём бота для отправки уведомлений
	baseBot, err := h.BotFactory.GetBaseBot(*integration.NotificationBotToken, h.Logger)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	bot, err := notification_bot.NewBot(baseBot, *integration.ErrorsNotificationChatID)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Отправляем уведомление
	bot.SendBugsnagNotification(string(bodyBytes))

	h.SendNewResponse(c, http.StatusOK, "OK")
}
