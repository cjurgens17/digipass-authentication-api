package handlers

import (
	"DigiPassAuthenticationApi/packages/models"
	"DigiPassAuthenticationApi/services"
	"DigiPassAuthenticationApi/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type AccountHandler struct{}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{}
}

type CreateAccountRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Email string `json:"email" validate:"required,email,max=255"`
}

func (h *AccountHandler) Create(c *echo.Context) error {
	var req CreateAccountRequest

	// Bind request body
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				switch e.Tag() {
				case "required":
					errors[field] = fmt.Sprintf("%s is required", field)
				case "email":
					errors[field] = fmt.Sprintf("%s must be a valid email address", field)
				case "min":
					errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
				case "max":
					errors[field] = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
				default:
					errors[field] = fmt.Sprintf("%s is invalid", field)
				}
			}

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":  "validation failed",
				"fields": errors,
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation error occurred",
		})
	}

	// Sanitize inputs
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Additional security validation
	if strings.ContainsAny(req.Name, "<>\"'&") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Name contains invalid characters",
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
		if errors.Is(err, services.ErrAccountAlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Account already exists with given email",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create account",
		})
	}

	return c.JSON(http.StatusCreated, account)
}
