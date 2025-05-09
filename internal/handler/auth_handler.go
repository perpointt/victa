package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"victa/internal/service"
	"victa/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.SendResponse(c, http.StatusUnauthorized, nil, err.Error())
		return
	}

	utils.SendResponse(c, http.StatusOK, gin.H{"token": token}, "Login successful")
}

func (h AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
		CompanyID *int64 `json:"company_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.CompanyID)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	utils.SendResponse(c, http.StatusCreated, user, "User registered successfully")
}
