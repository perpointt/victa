package webhook

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"victa/internal/bot/bot_common"
	"victa/internal/logger"
	"victa/internal/service"
	"victa/internal/webhook/webhook_common"
)

type GitlabWebhookHandler struct {
	*webhook_common.BaseWebhook
}

func NewGitlabWebhookHandler(
	factory *bot_common.BotFactory,
	logger logger.Logger,
	jwtSvc *service.JWTService,
) *GitlabWebhookHandler {
	base := webhook_common.NewBaseWebhook(factory, logger, jwtSvc)
	return &GitlabWebhookHandler{
		BaseWebhook: base,
	}
}

func (h *GitlabWebhookHandler) Handle(c *gin.Context) {
	_, err := h.Authorize(c)
	if err != nil {
		h.SendNewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	raw, err := c.GetRawData()
	if err != nil {
		h.SendNewResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GitLab webhook payload: %s", raw)

	h.SendNewResponse(c, http.StatusOK, "OK")
}
