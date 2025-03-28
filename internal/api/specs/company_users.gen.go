// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// RemoveUsersFromCompanyJSONBody defines parameters for RemoveUsersFromCompany.
type RemoveUsersFromCompanyJSONBody struct {
	UserIds []int `json:"user_ids"`
}

// AddUsersToCompanyJSONBody defines parameters for AddUsersToCompany.
type AddUsersToCompanyJSONBody struct {
	UserIds []int `json:"user_ids"`
}

// RemoveUsersFromCompanyJSONRequestBody defines body for RemoveUsersFromCompany for application/json ContentType.
type RemoveUsersFromCompanyJSONRequestBody RemoveUsersFromCompanyJSONBody

// AddUsersToCompanyJSONRequestBody defines body for AddUsersToCompany for application/json ContentType.
type AddUsersToCompanyJSONRequestBody AddUsersToCompanyJSONBody

// CompanyUsersServerInterface represents all server handlers.
type CompanyUsersServerInterface interface {
	// Удаление пользователей из компании
	// (DELETE /company-users/{id})
	RemoveUsersFromCompany(c *gin.Context, id int)
	// Получение информации о пользователе компании
	// (GET /company-users/{id})
	GetUsersInCompany(c *gin.Context, id int)
	// Добавление пользователей в компанию
	// (POST /company-users/{id})
	AddUsersToCompany(c *gin.Context, id int)
}

// CompanyUsersServerInterfaceWrapper converts contexts to parameters.
type CompanyUsersServerInterfaceWrapper struct {
	Handler            CompanyUsersServerInterface
	HandlerMiddlewares []CompanyUsersMiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type CompanyUsersMiddlewareFunc func(c *gin.Context)

// RemoveUsersFromCompany operation middleware
func (siw *CompanyUsersServerInterfaceWrapper) RemoveUsersFromCompany(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RemoveUsersFromCompany(c, id)
}

// GetUsersInCompany operation middleware
func (siw *CompanyUsersServerInterfaceWrapper) GetUsersInCompany(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUsersInCompany(c, id)
}

// AddUsersToCompany operation middleware
func (siw *CompanyUsersServerInterfaceWrapper) AddUsersToCompany(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id int

	err = runtime.BindStyledParameter("simple", false, "id", c.Param("id"), &id)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AddUsersToCompany(c, id)
}

// CompanyUsersGinServerOptions provides options for the Gin server.
type CompanyUsersGinServerOptions struct {
	BaseURL      string
	Middlewares  []CompanyUsersMiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterCompanyUsersHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterCompanyUsersHandlers(router gin.IRouter, si CompanyUsersServerInterface) {
	RegisterCompanyUsersHandlersWithOptions(router, si, CompanyUsersGinServerOptions{})
}

// RegisterCompanyUsersHandlersWithOptions creates http.Handler with additional options
func RegisterCompanyUsersHandlersWithOptions(router gin.IRouter, si CompanyUsersServerInterface, options CompanyUsersGinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := CompanyUsersServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.DELETE(options.BaseURL+"/company-users/:id", wrapper.RemoveUsersFromCompany)
	router.GET(options.BaseURL+"/company-users/:id", wrapper.GetUsersInCompany)
	router.POST(options.BaseURL+"/company-users/:id", wrapper.AddUsersToCompany)
}
