package router

import (
	"github.com/gin-gonic/gin"
	"victa/internal/handler"
)

// SetupRouter настраивает маршруты API версии v1.
func SetupRouter(companyHandler *handler.CompanyHandler, userHandler *handler.UserHandler /*, другие обработчики */) *gin.Engine {
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

		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Здесь можно добавить маршруты для приложений и аутентификации.
	}

	return r
}
