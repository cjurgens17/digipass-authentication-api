package routes

import (
	v1 "DigiPassAuthenticationApi/routes/v1"
	"github.com/labstack/echo/v5"
)

func SetUpRoutes(e *echo.Echo) {
	apiv1 := e.Group("v1")

	v1.RegisterAccountRoutes(apiv1)
}
