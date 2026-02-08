package v1

import (
	"DigiPassAuthenticationApi/handlers"
	"github.com/labstack/echo/v5"
)

func RegisterTenantRoutes(e *echo.Group) {
	v1Tenant := e.Group("/tenant")

	//Handler
	tenantHandler := handlers.NewAccountUsersHandler()

	println("Printing for placeholder right now", v1Tenant, tenantHandler)
}
