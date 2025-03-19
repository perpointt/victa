package utils

import "github.com/gin-gonic/gin"

// GetUserIDFromContext извлекает user_id из контекста, который устанавливается JWT-миддлварой.
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	uidFloat, ok := userIDInterface.(float64)
	if !ok {
		return 0, false
	}
	return int64(uidFloat), true
}
