package v1

import (
	"DigiPassAuthenticationApi/handlers"
	"github.com/labstack/echo/v5"
)

func RegisterAccountUsersRoutes(e *echo.Group) {
	v1AccountUsers := e.Group("/accountusers")

	//Handler
	accountUsersHandler := handlers.NewAccountUsersHandler()

	println("Printing for placeholder right now", v1AccountUsers, accountUsersHandler)
}
