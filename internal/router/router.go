package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	api "victa/internal/api/specs"
	"victa/internal/middleware"
)

// SetupRouter настраивает маршруты API версии v1.
// jwtSecret передаётся миддлваре для проверки токенов.
func SetupRouter(
	authHandler api.AuthServerInterface,
	companyHandler api.CompaniesServerInterface,
	companyUsersHandler api.CompanyUsersServerInterface,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	version := "v1"
	relativePath := fmt.Sprintf("/api/%s", version)

	group := r.Group(relativePath)

	api.RegisterAuthHandlers(group, authHandler)
	api.RegisterCompaniesHandlersWithOptions(group, companyHandler, api.CompaniesGinServerOptions{
		Middlewares: []api.CompaniesMiddlewareFunc{
			middleware.JWTAuthMiddleware(jwtSecret),
		},
	})
	api.RegisterCompanyUsersHandlersWithOptions(group, companyUsersHandler, api.CompanyUsersGinServerOptions{
		Middlewares: []api.CompanyUsersMiddlewareFunc{
			middleware.JWTAuthMiddleware(jwtSecret),
		},
	})

	return r
}
