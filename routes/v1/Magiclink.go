package v1

import (
	"DigiPassAuthenticationApi/handlers"
	"github.com/labstack/echo/v5"
)

func RegisterMagicLinkRoutes(e *echo.Group) {
	v1MagicLink := e.Group("/magiclink")

	// Handler
	magicLinkHandler := handlers.NewMagicLinkHandler()

	v1MagicLink.POST("/verify", magicLinkHandler.Verify)
	v1MagicLink.GET("/callback", magicLinkHandler.Callback)
}