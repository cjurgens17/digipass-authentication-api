package routes

import (
	"github.com/labstack/echo/v5"
	"DigiPassAuthenticationApi/services"
	"DigiPassAuthenticationApi/handlers"
	"gorm.io/gorm"
)

func SetUpRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/v1")
	
	accountService := services.NewAccountService(db)
	accountHandler := handlers.NewAccountHandler(accountService)

	api.POST("/account/new", accountHandler.Create)
}