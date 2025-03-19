package router

import (
	"github.com/gin-gonic/gin"
	"victa/internal/handler"
)

// SetupRouter настраивает маршруты API версии v1.
func SetupRouter(
	companyHandler *handler.CompanyHandler,
	userHandler *handler.UserHandler,
	appHandler *handler.AppHandler,
	authHandler *handler.AuthHandler,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		// Маршруты для аутентификации
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Маршруты для компаний
		companies := api.Group("/companies")
		{
			companies.POST("", companyHandler.CreateCompany)
			companies.GET("", companyHandler.GetCompanies)
			companies.GET("/:id", companyHandler.GetCompany)
			companies.PUT("/:id", companyHandler.UpdateCompany)
			companies.DELETE("/:id", companyHandler.DeleteCompany)
		}

		// Маршруты для пользователей
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Маршруты для приложений
		apps := api.Group("/apps")
		{
			apps.POST("", appHandler.CreateApp)
			apps.GET("", appHandler.GetApps)
			apps.GET("/:id", appHandler.GetApp)
			apps.PUT("/:id", appHandler.UpdateApp)
			apps.DELETE("/:id", appHandler.DeleteApp)
		}
	}

	return r
}
