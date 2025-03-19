package router

import (
	"github.com/gin-gonic/gin"
	"victa/internal/handler"
	"victa/internal/middleware"
)

// SetupRouter настраивает маршруты API версии v1.
// jwtSecret передаётся миддлваре для проверки токенов.
func SetupRouter(
	companyHandler *handler.CompanyHandler,
	userHandler *handler.UserHandler,
	appHandler *handler.AppHandler,
	authHandler *handler.AuthHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		// Открытые маршруты для аутентификации
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Защищенные маршруты: требуется валидный JWT
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware(jwtSecret))
		{
			companies := protected.Group("/companies")
			{
				companies.POST("", companyHandler.CreateCompany)
				companies.GET("", companyHandler.GetCompanies)
				companies.GET("/:id", companyHandler.GetCompany)
				companies.PUT("/:id", companyHandler.UpdateCompany)
				companies.DELETE("/:id", companyHandler.DeleteCompany)
			}

			// Эндпоинт для удаления пользователя из компании, вынесен в отдельную группу чтобы избежать конфликта с /companies/:id
			companyUsers := api.Group("/company-users")
			{
				companyUsers.GET("/:company_id", companyHandler.GetUsersInCompany)
				companyUsers.GET("/:company_id/:user_id", companyHandler.GetUserInCompany)
				companyUsers.DELETE("/:company_id/:user_id", companyHandler.DeleteUserFromCompany)
			}

			// Для эндпоинтов пользователя:
			users := api.Group("/user")
			{
				users.GET("/current", userHandler.GetCurrentUser)
				users.DELETE("/current", userHandler.DeleteAccount)
			}

			apps := protected.Group("/apps")
			{
				apps.POST("", appHandler.CreateApp)
				apps.GET("", appHandler.GetApps)
				apps.GET("/:id", appHandler.GetApp)
				apps.PUT("/:id", appHandler.UpdateApp)
				apps.DELETE("/:id", appHandler.DeleteApp)
			}
		}
	}

	return r
}
