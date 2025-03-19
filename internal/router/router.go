package router

import (
	"github.com/gin-gonic/gin"
	"victa/internal/handler"
)

// SetupRouter настраивает маршруты API версии v1.
func SetupRouter(companyHandler *handler.CompanyHandler /*, сюда можно добавить другие обработчики */) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		companies := api.Group("/companies")
		{
			companies.POST("", companyHandler.CreateCompany)
			companies.GET("", companyHandler.GetCompanies)
			companies.GET("/:id", companyHandler.GetCompany)
			companies.PUT("/:id", companyHandler.UpdateCompany)
			companies.DELETE("/:id", companyHandler.DeleteCompany)
		}
		// Добавьте маршруты для пользователей, приложений и аутентификации
	}

	return r
}
