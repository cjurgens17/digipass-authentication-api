package handlers

import (
	"DigiPassAuthenticationApi/packages/models"
	"DigiPassAuthenticationApi/services"
	"DigiPassAuthenticationApi/utils"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"net/http"
)

type AccountHandler struct{}

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

	var account *models.Account

	err := utils.WithTransaction(getDBFromContext(c), func(tx *gorm.DB) error {
		accountService := services.NewAccountService(tx)
		var err error
		account, err = accountService.CreateAccount(req.Name, req.Email)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err.Error() == "account already exits" {
			return c.JSON(http.StatusBadRequest, map[string]string {
				"error": "Account already exist with given email",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, account)
}
