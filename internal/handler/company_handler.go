package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"victa/internal/domain"
	"victa/internal/service"
)

// CompanyHandler обрабатывает HTTP запросы для компаний.
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, company)
}

// GetCompanies обрабатывает GET /api/v1/companies
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	companies, err := h.service.GetAllCompanies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, companies)
}

// GetCompany обрабатывает GET /api/v1/companies/:id
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	company, err := h.service.GetCompanyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, company)
}

// UpdateCompany обрабатывает PUT /api/v1/companies/:id
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var company domain.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	company.ID = id

	if err := h.service.UpdateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, company)
}

// DeleteCompany обрабатывает DELETE /api/v1/companies/:id
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteCompany(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
