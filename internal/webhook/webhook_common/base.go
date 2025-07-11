package webhook_common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"victa/internal/bot/bot_common"
	"victa/internal/domain"
	"victa/internal/logger"
	"victa/internal/service"
)

type BaseWebhook struct {
	BotFactory *bot_common.BotFactory
	Logger     logger.Logger
	jwtSvc     *service.JWTService
}

func NewBaseWebhook(
	botFactory *bot_common.BotFactory,
	logger logger.Logger,
	jwtSvc *service.JWTService,
) *BaseWebhook {
	return &BaseWebhook{
		BotFactory: botFactory,
		Logger:     logger,
		jwtSvc:     jwtSvc,
	}
}

func (wh *BaseWebhook) Authorize(c *gin.Context) (int64, error) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return 0, fmt.Errorf("missing authorization token")
	}
	return wh.jwtSvc.ParseToken(strings.TrimPrefix(auth, "Bearer "))
}

func (wh *BaseWebhook) SendNewResponse(
	c *gin.Context,
	code int,
	msg string,
) {
	c.AbortWithStatusJSON(code, domain.ApiResponse{Status: code, Message: msg})
}
