package utils

import (
	"github.com/gin-gonic/gin"
)

// SendResponse отправляет стандартизированный JSON ответ.
func SendResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, gin.H{
		"data":    data,
		"message": message,
		"status":  status,
	})
}

// AbortResponse отправляет стандартизированный JSON ответ.
func AbortResponse(c *gin.Context, status int, data interface{}, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"data":    data,
		"message": message,
		"status":  status,
	})
}
