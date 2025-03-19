package response

import (
	"github.com/gin-gonic/gin"
)

// APIResponse задаёт общий формат ответа.
type APIResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

// SendResponse отправляет JSON ответ в стандартизированном формате.
func SendResponse(c *gin.Context, status int, data interface{}, message string) {
	// Если data == nil и ожидается список, можно заменить на пустой срез.
	// Это можно расширить дополнительной логикой, если необходимо.
	c.JSON(status, APIResponse{
		Data:    data,
		Message: message,
		Status:  status,
	})
}
