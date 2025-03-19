package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"victa/internal/domain"
	"victa/internal/service"
)

// AppHandler обрабатывает HTTP-запросы для приложений.
type AppHandler struct {
	service service.AppService
}

// NewAppHandler создаёт новый AppHandler.
func NewAppHandler(service service.AppService) *AppHandler {
	return &AppHandler{service: service}
}

// CreateApp обрабатывает POST /api/v1/apps
func (h *AppHandler) CreateApp(c *gin.Context) {
	var app domain.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}
	if err := h.service.CreateApp(&app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data":    app,
		"message": "App created successfully",
		"status":  http.StatusCreated,
	})
}

// GetApps обрабатывает GET /api/v1/apps
func (h *AppHandler) GetApps(c *gin.Context) {
	apps, err := h.service.GetAllApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	// Если список равен nil, заменяем его на пустой срез
	if apps == nil {
		apps = []domain.App{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    apps,
		"message": "Apps retrieved successfully",
		"status":  http.StatusOK,
	})
}

// GetApp обрабатывает GET /api/v1/apps/:id
func (h *AppHandler) GetApp(c *gin.Context) {
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
	app, err := h.service.GetAppByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    app,
		"message": "App retrieved successfully",
		"status":  http.StatusOK,
	})
}

// UpdateApp обрабатывает PUT /api/v1/apps/:id
func (h *AppHandler) UpdateApp(c *gin.Context) {
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

	var app domain.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}
	app.ID = id

	if err := h.service.UpdateApp(&app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    app,
		"message": "App updated successfully",
		"status":  http.StatusOK,
	})
}

// DeleteApp обрабатывает DELETE /api/v1/apps/:id
func (h *AppHandler) DeleteApp(c *gin.Context) {
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
	if err := h.service.DeleteApp(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":    nil,
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"message": "App deleted successfully",
		"status":  http.StatusOK,
	})
}
