package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"victa/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// JWTAuthMiddleware проверяет наличие и валидность JWT токена в заголовке Authorization.
// При ошибке возвращает JSON-ответ в формате ApiResponse с кодом 401.
func JWTAuthMiddleware(jwtSecret string) func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.AbortResponse(c, http.StatusUnauthorized, nil, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.AbortResponse(c, http.StatusUnauthorized, nil, "Invalid Authorization header format")
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что используется HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.AbortResponse(c, http.StatusUnauthorized, nil, "Invalid or expired token")
			return
		}

		// Извлекаем claims и сохраняем полезные данные в контекст
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("company_id", claims["company_id"])
		}
		c.Next()
	}
}
