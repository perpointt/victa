package handler

import (
	"net/http"
	"strconv"
	"victa/internal/response"

	"github.com/gin-gonic/gin"
	"victa/internal/domain"
	"victa/internal/service"
)

// CompanyHandler обрабатывает HTTP-запросы для компаний.
type CompanyHandler struct {
	service service.CompanyService
}

// NewCompanyHandler создаёт новый CompanyHandler.
func NewCompanyHandler(service service.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

// CreateCompany обрабатывает POST /api/v1/companies
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company domain.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Проверяем наличие user_id в контексте (устанавливается JWT-миддлварой)
	userIDInterface, exists := c.Get("user_id")
	if exists {
		uidFloat, ok := userIDInterface.(float64)
		if !ok {
			response.SendResponse(c, http.StatusInternalServerError, nil, "Invalid user id type")
			return
		}
		userID := int64(uidFloat)
		if err := h.service.CreateCompanyAndLink(&company, userID); err != nil {
			response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
			return
		}
		response.SendResponse(c, http.StatusCreated, company, "Company created and linked successfully")
		return
	}

	// Если user_id не найден в контексте, создаем компанию без связи.
	if err := h.service.CreateCompany(&company); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusCreated, company, "Company created successfully")
}

// GetCompanies обрабатывает GET /api/v1/companies и возвращает компании, связанные с пользователем.
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	// Извлекаем user_id из контекста, который устанавливает JWT-миддлвара.
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	uidFloat, ok := userIDInterface.(float64)
	if !ok {
		response.SendResponse(c, http.StatusInternalServerError, nil, "Invalid user id type")
		return
	}
	userID := int64(uidFloat)

	companies, err := h.service.GetAllByUserID(userID)
	if err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	if companies == nil {
		companies = []domain.Company{}
	}

	response.SendResponse(c, http.StatusOK, companies, "Companies retrieved successfully")
}

// GetCompany обрабатывает GET /api/v1/companies/:id и возвращает компанию,
// если она связана с пользователем, чей id извлечён из JWT.
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	// Извлекаем user_id из контекста (устанавливается JWT-миддлварой).
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	uidFloat, ok := userIDInterface.(float64)
	if !ok {
		response.SendResponse(c, http.StatusInternalServerError, nil, "Invalid user id type")
		return
	}
	userID := int64(uidFloat)

	// Извлекаем company_id из параметров запроса.
	idStr := c.Param("id")
	companyID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
		return
	}

	// Получаем компанию, связанную с пользователем.
	company, err := h.service.GetCompanyByIDForUser(userID, companyID)
	if err != nil {
		response.SendResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, company, "Company retrieved successfully")
}

// UpdateCompany обрабатывает PUT /api/v1/companies/:id
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
		return
	}

	var company domain.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	company.ID = id

	if err := h.service.UpdateCompany(&company); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.SendResponse(c, http.StatusOK, company, "Company updated successfully")
}

// DeleteCompany обрабатывает DELETE /api/v1/companies/:id
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
		return
	}
	if err := h.service.DeleteCompany(id); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.SendResponse(c, http.StatusOK, nil, "Company deleted successfully")
}
