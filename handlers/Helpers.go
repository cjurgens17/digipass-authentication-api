package handlers

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func getDBFromContext(c *echo.Context) *gorm.DB {
	return c.Get("db").(*gorm.DB)
}