package handlers

import (
	"fmt"
	"net/http"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type MagiclinkHandler struct {}
type MagicLinkMetadata struct {
	ExpirationMinutes uint8 `json:"expirationMinutes" validate:"required,min=1,max=60"`
	EmailBody string `json:"emailBody" validate:"required,max=300"`
}
type MagicLinkVerify struct {
	ApiKey string `json:"apiKey" validate:"required,max=20"`
	EmailTo string `json:"emailTo" validate:"required,email"`
	EmailFrom string `json:"emailFrom" validate:"required,email"`
	RedirectUrl string `json:"redirectUrl" validate:"required,url"`
	Metadata MagicLinkMetadata `json:"metadata" validate:"required"`
}

func NewMagicLinkHandler() *MagiclinkHandler{
	return &MagiclinkHandler{}
}

func (h *MagiclinkHandler) Verify(c *echo.Context) error {
	var magicLinkVerify MagicLinkVerify
	//Bind Request Body
	if err := c.Bind(&magicLinkVerify); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	//Validate
	if err := c.Validate(&magicLinkVerify); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				switch e.Tag() {
				case "required":
					errors[field] = fmt.Sprintf("%s is required", field)
				case "email":
					errors[field] = fmt.Sprintf("%s must be a valid email address", field)
				case "url":
					errors[field] = fmt.Sprintf("%s must be a valid URL", field)
				case "min":
					errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
				case "max":
					errors[field] = fmt.Sprintf("%s must be at most %s characters", field, e.Param())
				default:
					errors[field] = fmt.Sprintf("%s is invalid", field)
				}
			}

			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "validation failed",
				"fields": errors,
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "validation error occurred",
		})
	}
	
	//temp return for now
	return nil
}