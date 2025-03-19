package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"victa/internal/domain"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data":    user,
		"message": "User created successfully",
		"status":  http.StatusCreated,
	})
}

// GetUsers обрабатывает GET /api/v1/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	// Если список равен nil, заменяем его на пустой срез
	if users == nil {
		users = []domain.User{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    users,
		"message": "Users retrieved successfully",
		"status":  http.StatusOK,
	})
}

// GetUser обрабатывает GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "Invalid id",
			"status":  http.StatusBadRequest,
		})
		return
	}
	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "User retrieved successfully",
		"status":  http.StatusOK,
	})
}

// UpdateUser обрабатывает PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "Invalid id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}
	user.ID = id

	if err := h.service.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "User updated successfully",
		"status":  http.StatusOK,
	})
}

// DeleteUser обрабатывает DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": "Invalid id",
			"status":  http.StatusBadRequest,
		})
		return
	}
	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "User deleted successfully",
		"status":  http.StatusOK,
	})
}
