package v1

import (
	"DigiPassAuthenticationApi/handlers"
	"github.com/labstack/echo/v5"
)

func RegisterAccountRoutes(e *echo.Group) {
	v1Account := e.Group("/account")

	//Handler
	accountHandler := handlers.NewAccountHandler()

	v1Account.POST("/account/new", accountHandler.Create)
}
