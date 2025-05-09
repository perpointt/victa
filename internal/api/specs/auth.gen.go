// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package api

import (
	"github.com/gin-gonic/gin"
)

// LoginJSONBody defines parameters for Login.
type LoginJSONBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterJSONBody defines parameters for Register.
type RegisterJSONBody struct {
	CompanyId *int   `json:"company_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginJSONBody

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody RegisterJSONBody

// AuthServerInterface represents all server handlers.
type AuthServerInterface interface {
	// Вход в систему
	// (POST /auth/login)
	Login(c *gin.Context)
	// Регистрация пользователя
	// (POST /auth/register)
	Register(c *gin.Context)
}

// AuthServerInterfaceWrapper converts contexts to parameters.
type AuthServerInterfaceWrapper struct {
	Handler            AuthServerInterface
	HandlerMiddlewares []AuthMiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type AuthMiddlewareFunc func(c *gin.Context)

// Login operation middleware
func (siw *AuthServerInterfaceWrapper) Login(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Login(c)
}

// Register operation middleware
func (siw *AuthServerInterfaceWrapper) Register(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Register(c)
}

// AuthGinServerOptions provides options for the Gin server.
type AuthGinServerOptions struct {
	BaseURL      string
	Middlewares  []AuthMiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterAuthHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterAuthHandlers(router gin.IRouter, si AuthServerInterface) {
	RegisterAuthHandlersWithOptions(router, si, AuthGinServerOptions{})
}

// RegisterAuthHandlersWithOptions creates http.Handler with additional options
func RegisterAuthHandlersWithOptions(router gin.IRouter, si AuthServerInterface, options AuthGinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := AuthServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/auth/login", wrapper.Login)
	router.POST(options.BaseURL+"/auth/register", wrapper.Register)
}
