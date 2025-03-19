package router

import (
	"github.com/gin-gonic/gin"
	"victa/internal/handler"
)

// SetupRouter настраивает маршруты API версии v1.
func SetupRouter(companyHandler *handler.CompanyHandler, userHandler *handler.UserHandler, appHandler *handler.AppHandler) *gin.Engine {
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

		apps := api.Group("/apps")
		{
			apps.POST("", appHandler.CreateApp)
			apps.GET("", appHandler.GetApps)
			apps.GET("/:id", appHandler.GetApp)
			apps.PUT("/:id", appHandler.UpdateApp)
			apps.DELETE("/:id", appHandler.DeleteApp)
		}
		// Добавьте другие маршруты, если потребуется
	}

	return r
}
