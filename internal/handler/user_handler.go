package handler

import (
	"net/http"
	"strconv"
	"victa/internal/utils"

	"github.com/gin-gonic/gin"
	"victa/internal/domain"
	"victa/internal/response"
	"victa/internal/service"
)

// UserHandler обрабатывает HTTP-запросы для пользователей.
type UserHandler struct {
	service service.UserService
}

// NewUserHandler создаёт новый UserHandler.
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUser обрабатывает GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
		return
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		response.SendResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, user, "User retrieved successfully")
}

// UpdateUser обрабатывает PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
		return
	}
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	user.ID = id
	if err := h.service.UpdateUser(&user); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, user, "User updated successfully")
}

// GetCurrentUser возвращает данные текущего пользователя.
// Эндпоинт: GET /user/current
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		response.SendResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, user, "User retrieved successfully")
}

// DeleteAccount удаляет аккаунт текущего пользователя.
// Эндпоинт: DELETE /user/current
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	if err := h.service.DeleteUser(userID); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, nil, "User account deleted successfully")
}
