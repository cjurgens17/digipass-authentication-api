package v1

import (
	"github.com/labstack/echo/v5"
	"DigiPassAuthenticationApi/handlers"
)

func RegisterAccountRoutes(e *echo.Group) {
	v1Account := e.Group("/account")

	//Handler
	accountHandler := handlers.NewAccountHandler()

	v1Account.POST("/account/new", accountHandler.Create)
}