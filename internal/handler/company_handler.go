package handler

//
//import (
//	"net/http"
//	"strconv"
//	"victa/internal/response"
//	"victa/internal/utils"
//
//	"github.com/gin-gonic/gin"
//	"victa/internal/domain"
//	"victa/internal/service"
//)
//
//// CompanyHandler обрабатывает HTTP-запросы для компаний.
//type CompanyHandler struct {
//	service            service.CompanyService
//	userService        service.UserService
//	userCompanyService service.UserCompanyService
//}
//
//// NewCompanyHandler создаёт новый CompanyHandler с зависимостями.
//func NewCompanyHandler(companyService service.CompanyService, userService service.UserService, userCompanyService service.UserCompanyService) *CompanyHandler {
//	return &CompanyHandler{
//		service:            companyService,
//		userService:        userService,
//		userCompanyService: userCompanyService,
//	}
//}
//
//// CreateCompany обрабатывает POST /api/v1/companies
//func (h *CompanyHandler) CreateCompany(c *gin.Context) {
//	var company domain.Company
//	if err := c.ShouldBindJSON(&company); err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
//		return
//	}
//
//	userID, ok := utils.GetUserIDFromContext(c)
//	if !ok {
//		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
//		return
//	}
//
//	if err := h.service.CreateCompanyAndLink(&company, userID); err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusCreated, company, "Company created and linked successfully")
//
//}
//
//// GetCompanies обрабатывает GET /api/v1/companies и возвращает компании, связанные с пользователем.
//func (h *CompanyHandler) GetCompanies(c *gin.Context) {
//	userID, ok := utils.GetUserIDFromContext(c)
//	if !ok {
//		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
//		return
//	}
//
//	companies, err := h.service.GetAllByUserID(userID)
//	if err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	if companies == nil {
//		companies = []domain.Company{}
//	}
//
//	response.SendResponse(c, http.StatusOK, companies, "Companies retrieved successfully")
//}
//
//// GetCompany обрабатывает GET /api/v1/companies/:id и возвращает компанию, если она связана с пользователем.
//func (h *CompanyHandler) GetCompany(c *gin.Context) {
//	userID, ok := utils.GetUserIDFromContext(c)
//	if !ok {
//		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
//		return
//	}
//
//	// Извлекаем company_id из URL.
//	idStr := c.Param("id")
//	companyID, err := strconv.ParseInt(idStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
//		return
//	}
//
//	// Проверяем, что текущий пользователь имеет доступ к этой компании.
//	company, err := h.service.GetCompanyByIDForUser(userID, companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusNotFound, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusOK, company, "Company retrieved successfully")
//}
//
//// GetUsersInCompany обрабатывает GET /api/v1/companies/:id/users и возвращает список пользователей, связанных с этой компанией.
//func (h *CompanyHandler) GetUsersInCompany(c *gin.Context) {
//	userID, ok := utils.GetUserIDFromContext(c)
//	if !ok {
//		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
//		return
//	}
//
//	// Извлекаем company_id из параметров URL.
//	companyIDStr := c.Param("company_id")
//	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
//		return
//	}
//
//	// Проверяем, что текущий пользователь имеет доступ к данной компании.
//	_, err = h.service.GetCompanyByIDForUser(userID, companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusForbidden, nil, "Access denied or company not found")
//		return
//	}
//
//	// Получаем список пользователей в компании.
//	users, err := h.userService.GetAllUsersByCompany(companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	if users == nil {
//		users = []domain.User{}
//	}
//	response.SendResponse(c, http.StatusOK, users, "Users retrieved successfully")
//}
//
//// DeleteCompany обрабатывает DELETE /companies/:id.
//// Компания удаляется, если текущий пользователь имеет роль "admin" в ней.
//func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
//	userID, ok := utils.GetUserIDFromContext(c)
//	if !ok {
//		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
//		return
//	}
//	companyIDStr := c.Param("id")
//	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
//		return
//	}
//
//	isAdmin, err := h.userCompanyService.IsAdmin(userID, companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
//		return
//	}
//	if !isAdmin {
//		response.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
//		return
//	}
//
//	if err := h.service.DeleteCompany(companyID); err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusOK, nil, "Company deleted successfully")
//}
//
//// UpdateCompany обрабатывает PUT /companies/:id.
//// Обновление разрешено только, если пользователь является администратором компании.
//func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
//	userID, ok := utils.GetUserIDFromContext(c)
//	if !ok {
//		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
//		return
//	}
//	companyIDStr := c.Param("id")
//	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
//		return
//	}
//
//	isAdmin, err := h.userCompanyService.IsAdmin(userID, companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
//		return
//	}
//	if !isAdmin {
//		response.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
//		return
//	}
//
//	var company domain.Company
//	if err := c.ShouldBindJSON(&company); err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, err.Error())
//		return
//	}
//	company.ID = companyID
//
//	if err := h.service.UpdateCompany(&company); err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusOK, company, "Company updated successfully")
//}
//
//// DeleteUserFromCompany удаляет пользователя из компании, если вызывающий пользователь является администратором.
//// Эндпоинт: DELETE /company-users/:company_id/:user_id
//func (h *CompanyHandler) DeleteUserFromCompany(c *gin.Context) {
//	adminID, ok := utils.GetUserIDFromContext(c)
//	if !ok {
//		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
//		return
//	}
//
//	companyIDStr := c.Param("company_id")
//	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
//		return
//	}
//	targetUserIDStr := c.Param("user_id")
//	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid user id")
//		return
//	}
//
//	// Проверяем, что вызывающий пользователь является администратором компании.
//	isAdmin, err := h.userCompanyService.IsAdmin(adminID, companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusNotFound, nil, "Company not found or access denied")
//		return
//	}
//	if !isAdmin {
//		response.SendResponse(c, http.StatusForbidden, nil, "Access denied: insufficient permissions")
//		return
//	}
//
//	// Удаляем связь пользователя с компанией.
//	err = h.userCompanyService.RemoveUserFromCompany(targetUserID, companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusInternalServerError, nil, err.Error())
//		return
//	}
//	response.SendResponse(c, http.StatusOK, nil, "User removed from company successfully")
//}
//
//func (h *CompanyHandler) GetUserInCompany(c *gin.Context) {
//	// Извлекаем ID вызывающего пользователя из контекста.
//	requesterID, ok := utils.GetUserIDFromContext(c)
//	if !ok {
//		response.SendResponse(c, http.StatusUnauthorized, nil, "User not authenticated")
//		return
//	}
//
//	// Извлекаем company_id и user_id из URL.
//	companyIDStr := c.Param("company_id")
//	companyID, err := strconv.ParseInt(companyIDStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid company id")
//		return
//	}
//	targetUserIDStr := c.Param("user_id")
//	targetUserID, err := strconv.ParseInt(targetUserIDStr, 10, 64)
//	if err != nil {
//		response.SendResponse(c, http.StatusBadRequest, nil, "Invalid user id")
//		return
//	}
//
//	// Проверяем, что вызывающий пользователь является участником компании.
//	_, err = h.userCompanyService.GetUserRole(requesterID, companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusForbidden, nil, "Access denied: you are not a member of this company")
//		return
//	}
//
//	// Проверяем, что целевой пользователь состоит в компании.
//	_, err = h.userCompanyService.GetUserRole(targetUserID, companyID)
//	if err != nil {
//		response.SendResponse(c, http.StatusNotFound, nil, "Target user not found in company")
//		return
//	}
//
//	// Получаем полную информацию о целевом пользователе.
//	user, err := h.userService.GetUserByID(targetUserID)
//	if err != nil {
//		response.SendResponse(c, http.StatusNotFound, nil, "User not found")
//		return
//	}
//	response.SendResponse(c, http.StatusOK, user, "User retrieved successfully")
//}
