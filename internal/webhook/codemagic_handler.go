package webhook

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
	"victa/internal/bot/notification_bot"
	"victa/internal/domain"
	appErr "victa/internal/errors"
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

	apiKey, err := h.companySvc.GetCompanySecret(ctx, companyID, domain.SecretCodemagicApiKey)
	if err != nil && !errors.Is(err, appErr.ErrSecretNotFound) {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if apiKey == nil {
		h.SendNewResponse(c, http.StatusBadRequest, "Codemagic API key not found")
		return
	}

	cmCtx, cancelCM := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCM()

	build, err := h.codemagicSvc.GetBuildByID(cmCtx, payload.BuildID, string(apiKey))
	if err != nil {
		h.SendNewResponse(c, http.StatusBadGateway, err.Error())
		return
	}

	for i, art := range build.Build.Artefacts {
		if strings.EqualFold(art.Type, "apk") {
			url, err := h.codemagicSvc.GetArtifactPublicURL(cmCtx, art.Path, string(apiKey))
			if err != nil {
				h.Logger.Warn("codemagic public URL: %v", err)
				break
			}
			build.Build.Artefacts[i].PublicURL = url
			break
		}
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

	chatID, err := h.companySvc.GetCompanySecret(ctx, companyID, domain.SecretDeployNotificationChatID)
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

	bot.SendDeployNotification(build.Application, build.Build)

	h.SendNewResponse(c, http.StatusOK, "OK")
}
