package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"victa/internal/service"
	"victa/internal/utils"
)

// CompanyUsersHandler обрабатывает HTTP-запросы для пользователей.
type CompanyUsersHandler struct {
	companyUserService service.CompanyUserService
	userService        service.UserService
}

// NewCompanyUsersHandler создаёт новый UserHandler.
func NewCompanyUsersHandler(companyUserService service.CompanyUserService, userService service.UserService) *CompanyUsersHandler {
	return &CompanyUsersHandler{
		companyUserService: companyUserService,
		userService:        userService,
	}
}

func (h CompanyUsersHandler) GetUsersInCompany(c *gin.Context, id int) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	isUserInCompany, err := h.companyUserService.IsUserInCompany(userID, int64(id))
	if err != nil {
		utils.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
		return
	}

	if !isUserInCompany {
		utils.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
		return
	}

	users, err := h.userService.GetUsersInCompany(int64(id))
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.SendResponse(c, http.StatusOK, users, "Company users retrieved successfully")
}

func (h CompanyUsersHandler) AddUsersToCompany(c *gin.Context, id int) {
	// Проверяем, аутентифицирован ли пользователь
	currentUserID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	// Проверяем, что текущий пользователь является администратором компании
	isAdmin, err := h.companyUserService.IsUserAdminInCompany(currentUserID, int64(id))
	if err != nil {
		utils.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
		return
	}
	if !isAdmin {
		utils.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
		return
	}

	// Парсим тело запроса
	var req struct {
		UserIDs []int64 `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Вызываем сервис для добавления пользователей в компанию
	if err := h.companyUserService.LinkUsersWithCompany(req.UserIDs, int64(id)); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.SendResponse(c, http.StatusOK, nil, "Users added to company successfully")
}

func (h CompanyUsersHandler) RemoveUsersFromCompany(c *gin.Context, id int) {
	// Проверяем, аутентифицирован ли пользователь
	currentUserID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	// Проверяем, что текущий пользователь является администратором компании
	isAdmin, err := h.companyUserService.IsUserAdminInCompany(currentUserID, int64(id))
	if err != nil {
		utils.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
		return
	}
	if !isAdmin {
		utils.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
		return
	}

	// Парсим тело запроса
	var req struct {
		UserIDs []int64 `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Вызываем сервис для удаления пользователей из компании
	if err := h.companyUserService.UnlinkUsersCompany(req.UserIDs, int64(id)); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.SendResponse(c, http.StatusOK, nil, "Users removed from company successfully")
}
