package handlers

import (
	"DigiPassAuthenticationApi/services"
	"net/http"
	"github.com/labstack/echo/v5"
)

type AccountHandler struct {
	service *services.AccountService
}

func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
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

	account, err := h.service.CreateAccount(req.Name, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, account)
}
