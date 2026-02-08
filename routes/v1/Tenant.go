package v1

import (
	"github.com/labstack/echo/v5"
	"DigiPassAuthenticationApi/handlers"
)

func RegisterTenantRoutes(e *echo.Group) {
	v1Tenant := e.Group("/tenant")

	//Handler
	tenantHandler := handlers.NewAccountUsersHandler()

	println("Printing for placeholder right now",v1Tenant, tenantHandler)
}