package handler

import (
	"github.com/gin-gonic/gin"
	"victa/internal/service"
)

// CompanyUsersHandler обрабатывает HTTP-запросы для пользователей.
type CompanyUsersHandler struct {
	service service.CompanyUsersService
}

// NewCompanyUsersHandler создаёт новый UserHandler.
func NewCompanyUsersHandler(service service.CompanyUsersService) *CompanyUsersHandler {
	return &CompanyUsersHandler{service: service}
}

func (c2 CompanyUsersHandler) GetCompanyUsers(c *gin.Context, id int) {
	//TODO implement me
	panic("implement me")
}

func (c2 CompanyUsersHandler) AddCompanyUsers(c *gin.Context, id int) {
	//TODO implement me
	panic("implement me")
}

func (c2 CompanyUsersHandler) RemoveCompanyUsers(c *gin.Context, id int) {
	//TODO implement me
	panic("implement me")
}
