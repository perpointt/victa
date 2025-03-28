package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"victa/internal/domain"
	"victa/internal/service"
	"victa/internal/utils"
)

type CompanyHandler struct {
	service            service.CompanyService
	userService        service.UserService
	userCompanyService service.CompanyUserService
}

func NewCompanyHandler(companyService service.CompanyService, userService service.UserService, userCompanyService service.CompanyUserService) *CompanyHandler {
	return &CompanyHandler{
		service:            companyService,
		userService:        userService,
		userCompanyService: userCompanyService,
	}
}

// GetCompanies извлекает user_id из JWT и возвращает список компаний для данного пользователя.
func (h CompanyHandler) GetCompanies(c *gin.Context) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	companies, err := h.service.GetCompanies(userID)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	if companies == nil {
		companies = []domain.Company{}
	}

	utils.SendResponse(c, http.StatusOK, companies, "Companies retrieved successfully")
}

func (h CompanyHandler) CreateCompany(c *gin.Context) {
	var company domain.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	if err := h.service.CreateCompany(&company, userID); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.SendResponse(c, http.StatusCreated, company, "Company created successfully")
}

func (h CompanyHandler) GetCompany(c *gin.Context, id int) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	// Проверяем, что текущий пользователь имеет доступ к этой компании.
	company, err := h.service.GetCompanyByID(userID, int64(id))
	if err != nil {
		utils.SendResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}
	utils.SendResponse(c, http.StatusOK, company, "Company retrieved successfully")
}

func (h CompanyHandler) UpdateCompany(c *gin.Context, id int) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	isAdmin, err := h.userCompanyService.IsUserAdminInCompany(userID, int64(id))
	if err != nil {
		utils.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
		return
	}
	if !isAdmin {
		utils.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
		return
	}

	var company domain.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	company.ID = int64(id)

	if err := h.service.UpdateCompany(&company); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	utils.SendResponse(c, http.StatusOK, company, "Company updated successfully")
}

func (h CompanyHandler) DeleteCompany(c *gin.Context, id int) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		utils.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	isAdmin, err := h.userCompanyService.IsUserAdminInCompany(userID, int64(id))
	if err != nil {
		utils.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
		return
	}
	if !isAdmin {
		utils.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
		return
	}

	if err := h.service.DeleteCompany(int64(id)); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	utils.SendResponse(c, http.StatusOK, nil, "Company deleted successfully")
}
