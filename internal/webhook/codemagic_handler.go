package webhook

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"victa/internal/bot/notification_bot"
	"victa/internal/logger"

	"victa/internal/bot/bot_common"
	"victa/internal/service"
)

type apiResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CodemagicWebhookHandler struct {
	factory      *bot_common.BotFactory
	logger       logger.Logger
	jwtSvc       *service.JWTService
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
	return &CodemagicWebhookHandler{
		factory,
		logger,
		jwtSvc,
		companySvc,
		codemagicSvc,
	}
}

func (h *CodemagicWebhookHandler) Handle(c *gin.Context) {
	companyID, err := h.authorize(c)
	if err != nil {
		h.respondJSON(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	var payload struct {
		BuildID string `json:"build_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.respondJSON(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	integration, err := h.companySvc.GetCompanyIntegrationByID(companyID)
	if err != nil {
		h.respondJSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if integration == nil {
		h.respondJSON(c, http.StatusBadRequest, "company not found", nil)
		return
	}

	buildResp, err := h.codemagicSvc.GetBuildByID(payload.BuildID, *integration.CodemagicAPIKey)
	if err != nil {
		h.respondJSON(c, http.StatusBadGateway, err.Error(), nil)
		return
	}
	for idx, art := range buildResp.Build.Artefacts {
		if strings.EqualFold(art.Type, "apk") {
			publicURL, err := h.codemagicSvc.GetArtifactPublicURL(
				art.Path, *integration.CodemagicAPIKey,
			)
			if err != nil {
				h.logger.Errorf(err.Error())
				break
			}
			buildResp.Build.Artefacts[idx].PublicURL = publicURL
			break
		}
	}

	baseBot, err := h.factory.GetBaseBot(*integration.NotificationBotToken, h.logger)
	if err != nil {
		h.respondJSON(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	bot, err := notification_bot.NewBot(baseBot, *integration.NotificationChatID)
	if err != nil {
		h.respondJSON(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	bot.SendNewNotification(buildResp.Application, buildResp.Build)

	h.respondJSON(c, http.StatusOK, "OK", nil)
}

func (h *CodemagicWebhookHandler) authorize(c *gin.Context) (int64, error) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return 0, fmt.Errorf("missing bearer token")
	}
	return h.jwtSvc.ParseToken(strings.TrimPrefix(auth, "Bearer "))
}

func (h *CodemagicWebhookHandler) respondJSON(
	c *gin.Context,
	code int,
	msg string,
	headers map[string]string,
) {
	for k, v := range headers {
		c.Header(k, v)
	}
	c.AbortWithStatusJSON(code, apiResponse{Status: code, Message: msg})
}
