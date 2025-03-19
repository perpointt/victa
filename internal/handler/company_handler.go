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
	service            service.CompanyService
	userService        service.UserService
	userCompanyService service.UserCompanyService
}

// NewCompanyHandler создаёт новый CompanyHandler с зависимостями.
func NewCompanyHandler(companyService service.CompanyService, userService service.UserService, userCompanyService service.UserCompanyService) *CompanyHandler {
	return &CompanyHandler{
		service:            companyService,
		userService:        userService,
		userCompanyService: userCompanyService,
	}
}

// helper: извлечение userID из контекста
func getUserIDFromContext(c *gin.Context) (int64, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	uidFloat, ok := userIDInterface.(float64)
	if !ok {
		return 0, false
	}
	return int64(uidFloat), true
}

// CreateCompany обрабатывает POST /api/v1/companies
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company domain.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	if err := h.service.CreateCompanyAndLink(&company, userID); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusCreated, company, "Company created and linked successfully")

}

// GetCompanies обрабатывает GET /api/v1/companies и возвращает компании, связанные с пользователем.
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

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

// GetCompany обрабатывает GET /api/v1/companies/:id и возвращает компанию, если она связана с пользователем.
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	// Извлекаем company_id из URL.
	idStr := c.Param("id")
	companyID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
		return
	}

	// Проверяем, что текущий пользователь имеет доступ к этой компании.
	company, err := h.service.GetCompanyByIDForUser(userID, companyID)
	if err != nil {
		response.SendResponse(c, http.StatusNotFound, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, company, "Company retrieved successfully")
}

// GetUsersInCompany обрабатывает GET /api/v1/companies/:id/users и возвращает список пользователей, связанных с этой компанией.
func (h *CompanyHandler) GetUsersInCompany(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}

	// Извлекаем company_id из параметров URL.
	companyIDStr := c.Param("id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
		return
	}

	// Проверяем, что текущий пользователь имеет доступ к данной компании.
	_, err = h.service.GetCompanyByIDForUser(userID, companyID)
	if err != nil {
		response.SendResponse(c, http.StatusForbidden, nil, "Access denied or company not found")
		return
	}

	// Получаем список пользователей в компании.
	users, err := h.userService.GetAllUsersByCompany(companyID)
	if err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	if users == nil {
		users = []domain.User{}
	}
	response.SendResponse(c, http.StatusOK, users, "Users retrieved successfully")
}

// DeleteCompany обрабатывает DELETE /companies/:id.
// Компания удаляется, если текущий пользователь имеет роль "admin" в ней.
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	companyIDStr := c.Param("id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
		return
	}

	isAdmin, err := h.userCompanyService.IsAdmin(userID, companyID)
	if err != nil {
		response.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
		return
	}
	if !isAdmin {
		response.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
		return
	}

	if err := h.service.DeleteCompany(companyID); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, nil, "Company deleted successfully")
}

// UpdateCompany обрабатывает PUT /companies/:id.
// Обновление разрешено только, если пользователь является администратором компании.
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
		return
	}
	companyIDStr := c.Param("id")
	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
	if err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
		return
	}

	isAdmin, err := h.userCompanyService.IsAdmin(userID, companyID)
	if err != nil {
		response.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
		return
	}
	if !isAdmin {
		response.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
		return
	}

	var company domain.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	company.ID = companyID

	if err := h.service.UpdateCompany(&company); err != nil {
		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.SendResponse(c, http.StatusOK, company, "Company updated successfully")
}
