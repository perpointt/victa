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
	if err := h.service.CreateCompany(&company); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	response.SendResponse(c, http.StatusCreated, company, "Company created successfully")
}

// GetCompanies обрабатывает GET /api/v1/companies
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	companies, err := h.service.GetAllCompanies()
	if err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	if companies == nil {
		companies = []domain.Company{}
	}
	response.SendResponse(c, http.StatusOK, companies, "Companies retrieved successfully")
}

// GetCompany обрабатывает GET /api/v1/companies/:id
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid id")
		return
	}
	company, err := h.service.GetCompanyByID(id)
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
