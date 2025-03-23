package handler

//
//import (
//	"net/http"
//	"strconv"
//
//	"github.com/gin-gonic/gin"
//	"victa/internal/domain"
//	"victa/internal/response"
//	"victa/internal/service"
//)
//
//// AppHandler обрабатывает HTTP-запросы для приложений.
//type AppHandler struct {
//	service service.AppService
//}
//
//// NewAppHandler создаёт новый AppHandler.
//func NewAppHandler(service service.AppService) *AppHandler {
//	return &AppHandler{service: service}
//}
//
//// CreateApp обрабатывает POST /api/v1/apps
//func (h *AppHandler) CreateApp(c *gin.Context) {
//	var app domain.App
//	if err := c.ShouldBindJSON(&app); err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
//		return
//	}
//	if err := h.service.CreateApp(&app); err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusCreated, app, "App created successfully")
//}
//
//// GetApps обрабатывает GET /api/v1/apps
//func (h *AppHandler) GetApps(c *gin.Context) {
//	apps, err := h.service.GetAllApps()
//	if err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	if apps == nil {
//		apps = []domain.App{}
//	}
//	response.SendResponse(c, http.StatusOK, apps, "Apps retrieved successfully")
//}
//
//// GetApp обрабатывает GET /api/v1/apps/:id
//func (h *AppHandler) GetApp(c *gin.Context) {
//	idStr := c.Param("id")
//	id, err := strconv.ParseInt(idStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
//		return
//	}
//	app, err := h.service.GetAppByID(id)
//	if err != nil {
//		response.SendResponse(c, http.StatusNotFound, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusOK, app, "App retrieved successfully")
//}
//
//// UpdateApp обрабатывает PUT /api/v1/apps/:id
//func (h *AppHandler) UpdateApp(c *gin.Context) {
//	idStr := c.Param("id")
//	id, err := strconv.ParseInt(idStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
//		return
//	}
//	var app domain.App
//	if err := c.ShouldBindJSON(&app); err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
//		return
//	}
//	app.ID = id
//	if err := h.service.UpdateApp(&app); err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusOK, app, "App updated successfully")
//}
//
//// DeleteApp обрабатывает DELETE /api/v1/apps/:id
//func (h *AppHandler) DeleteApp(c *gin.Context) {
//	idStr := c.Param("id")
//	id, err := strconv.ParseInt(idStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
//		return
//	}
//	if err := h.service.DeleteApp(id); err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusOK, nil, "App deleted successfully")
//}
