package handlers

import (
	"DigiPassAuthenticationApi/services"
	"net/http"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type AccountHandler struct {}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{}
}

func (h *AccountHandler) Create(c *echo.Context) error {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	accountService := services.NewAccountService(c.Get("db").(*gorm.DB))

	account, err := accountService.CreateAccount(req.Name, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, account)
}
