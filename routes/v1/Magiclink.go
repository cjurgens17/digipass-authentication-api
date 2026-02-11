package v1

import (
	"DigiPassAuthenticationApi/handlers"
	"github.com/labstack/echo/v5"
)

func RegisterMagicLinkRoutes(e *echo.Group) {
	v1Account := e.Group("/magiclink")

	//Handler
	accountHandler := handlers.NewAccountHandler()

	v1Account.POST("/new", accountHandler.Create)
}