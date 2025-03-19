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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateApp(&app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, app)
}

// GetApps обрабатывает GET /api/v1/apps
func (h *AppHandler) GetApps(c *gin.Context) {
	apps, err := h.service.GetAllApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apps)
}

// GetApp обрабатывает GET /api/v1/apps/:id
func (h *AppHandler) GetApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	app, err := h.service.GetAppByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, app)
}

// UpdateApp обрабатывает PUT /api/v1/apps/:id
func (h *AppHandler) UpdateApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var app domain.App
	if err := c.ShouldBindJSON(&app); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	app.ID = id

	if err := h.service.UpdateApp(&app); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, app)
}

// DeleteApp обрабатывает DELETE /api/v1/apps/:id
func (h *AppHandler) DeleteApp(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteApp(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
