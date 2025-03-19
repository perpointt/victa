package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"victa/internal/service"
)

// AuthHandler обрабатывает запросы для аутентификации.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler создаёт новый AuthHandler.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register обрабатывает POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
		CompanyID *int64 `json:"company_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.CompanyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    user,
		"message": "User registered successfully",
		"status":  http.StatusCreated,
	})
}

// Login обрабатывает POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"token": token},
		"message": "Login successful",
		"status":  http.StatusOK,
	})
}
