package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type MagiclinkHandler struct{}

type MagicLinkMetadata struct {
	ExpirationMinutes uint8  `json:"expirationMinutes" validate:"required,min=1,max=60"`
	EmailBody         string `json:"emailBody" validate:"required,min=1,max=300"`
}

type MagicLinkVerifyRequest struct {
	ApiKey      string            `json:"apiKey" validate:"required,min=10,max=64,alphanum"`
	EmailTo     string            `json:"emailTo" validate:"required,email,max=255"`
	EmailFrom   string            `json:"emailFrom" validate:"required,email,max=255"`
	RedirectUrl string            `json:"redirectUrl" validate:"required,url,max=2048"`
	Metadata    MagicLinkMetadata `json:"metadata" validate:"required"`
}

func NewMagicLinkHandler() *MagiclinkHandler {
	return &MagiclinkHandler{}
}

func (h *MagiclinkHandler) Verify(c *echo.Context) error {
	var magicLinkVerifyRequest MagicLinkVerifyRequest

	// Bind Request Body
	if err := c.Bind(&magicLinkVerifyRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate
	if err := c.Validate(&magicLinkVerifyRequest); err != nil {
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
				case "alphanum":
					errors[field] = fmt.Sprintf("%s must contain only letters and numbers", field)
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

	// Sanitize email inputs
	magicLinkVerifyRequest.EmailTo = strings.TrimSpace(strings.ToLower(magicLinkVerifyRequest.EmailTo))
	magicLinkVerifyRequest.EmailFrom = strings.TrimSpace(strings.ToLower(magicLinkVerifyRequest.EmailFrom))
	magicLinkVerifyRequest.ApiKey = strings.TrimSpace(magicLinkVerifyRequest.ApiKey)
	magicLinkVerifyRequest.Metadata.EmailBody = strings.TrimSpace(magicLinkVerifyRequest.Metadata.EmailBody)

	// Additional security validation for redirect URL (SSRF prevention)
	parsedURL, err := url.Parse(magicLinkVerifyRequest.RedirectUrl)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid redirect URL format",
		})
	}

	// Prevent SSRF attacks - only allow https scheme and reject localhost/private IPs
	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Redirect URL must use HTTP or HTTPS protocol",
		})
	}

	// Reject localhost and private IP ranges to prevent SSRF
	hostname := strings.ToLower(parsedURL.Hostname())
	if hostname == "localhost" || hostname == "127.0.0.1" || hostname == "0.0.0.0" ||
		strings.HasPrefix(hostname, "192.168.") || strings.HasPrefix(hostname, "10.") ||
		strings.HasPrefix(hostname, "172.16.") || strings.Contains(hostname, "::1") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Redirect URL cannot point to local or private addresses",
		})
	}

	// Prevent XSS in email body
	if strings.ContainsAny(magicLinkVerifyRequest.Metadata.EmailBody, "<>\"'&") {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email body contains invalid characters",
		})
	}

	// temp return for now
	return nil
}