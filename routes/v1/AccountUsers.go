package v1

import (
	"github.com/labstack/echo/v5"
	"DigiPassAuthenticationApi/handlers"
)

func RegisterAccountUsersRoutes(e *echo.Group) {
	v1AccountUsers := e.Group("/accountusers")

	//Handler
	accountUsersHandler := handlers.NewAccountUsersHandler()

	println("Printing for placeholder right now",v1AccountUsers, accountUsersHandler)
}