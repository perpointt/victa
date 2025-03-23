package handler

import (
	"github.com/gin-gonic/gin"
	"victa/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func (a AuthHandler) Login(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a AuthHandler) Register(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

//
//// Register обрабатывает POST /api/v1/auth/register
//func (h *AuthHandler) Register(c *gin.Context) {
//	var req struct {
//		Email     string `json:"email" binding:"required,email"`
//		Password  string `json:"password" binding:"required"`
//		CompanyID *int64 `json:"company_id"`
//	}
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
//		return
//	}
//
//	user, err := h.authService.Register(req.Email, req.Password, req.CompanyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
//		return
//	}
//
//	response.SendResponse(c, http.StatusCreated, user, "User registered successfully")
//}
//
//// Login обрабатывает POST /api/v1/auth/login
//func (h *AuthHandler) Login(c *gin.Context) {
//	var req struct {
//		Email    string `json:"email" binding:"required,email"`
//		Password string `json:"password" binding:"required"`
//	}
//
//	if err := c.ShouldBindJSON(&req); err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
//		return
//	}
//
//	token, err := h.authService.Login(req.Email, req.Password)
//	if err != nil {
//		response.SendResponse(c, http.StatusUnauthorized, nil, err.Error())
//		return
//	}
//
//	response.SendResponse(c, http.StatusOK, gin.H{"token": token}, "Login successful")
//}
