package handler

import (
	"net/http"
	"strconv"

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

// CreateUser обрабатывает POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if err := h.service.CreateUser(&user); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusCreated, user, "User created successfully")
}

// GetUsers обрабатывает GET /api/v1/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	if users == nil {
		users = []domain.User{}
	}
	response.SendResponse(c, http.StatusOK, users, "Users retrieved successfully")
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

// DeleteUser обрабатывает DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
		return
	}
	if err := h.service.DeleteUser(id); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, nil, "User deleted successfully")
}
