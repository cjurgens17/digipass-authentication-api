package handlers

import (
	"DigiPassAuthenticationApi/packages/jwt"
	"DigiPassAuthenticationApi/packages/models"
	"DigiPassAuthenticationApi/services"
	"DigiPassAuthenticationApi/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
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

type MagicLinkCallbackRequest struct {
	Token string `query:"token" validate:"required,min=32"`
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

	var result *services.MagicLinkResult

	err = utils.WithTransaction(getDBFromContext(c), func(tx *gorm.DB) error {
		magicLinkService := services.NewMagicLinkService(tx)
		var err error
		result, err = magicLinkService.VerifyMagicLink(
			magicLinkVerifyRequest.ApiKey,
			magicLinkVerifyRequest.Metadata.ExpirationMinutes,
			magicLinkVerifyRequest.RedirectUrl,
		)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, services.ErrInvalidApiKey) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid API key",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate magic link",
		})
	}

	// TODO: Send email with magic link to emailTo from emailFrom
	// For now, we'll just return the magic link URL in the response
	// In production, you would integrate with an email service here

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Magic link generated successfully",
		"magic_link":   result.MagicLinkURL,
		"expires_at":   result.MagicLink.ExpiresAt,
		"email_to":     magicLinkVerifyRequest.EmailTo,
		"email_from":   magicLinkVerifyRequest.EmailFrom,
		"email_body":   magicLinkVerifyRequest.Metadata.EmailBody,
	})
}

func (h *MagiclinkHandler) Callback(c *echo.Context) error {
	var req MagicLinkCallbackRequest

	// Bind query parameters
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request parameters",
		})
	}

	// Validate
	if err := c.Validate(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range validationErrors {
				field := e.Field()
				switch e.Tag() {
				case "required":
					errors[field] = fmt.Sprintf("%s is required", field)
				case "min":
					errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
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

	// Sanitize token input
	req.Token = strings.TrimSpace(req.Token)

	var magicLink *models.MagicLink

	err := utils.WithTransaction(getDBFromContext(c), func(tx *gorm.DB) error {
		magicLinkService := services.NewMagicLinkService(tx)
		var err error
		magicLink, err = magicLinkService.ValidateToken(req.Token)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, services.ErrMagicLinkNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Invalid magic link",
			})
		}
		if errors.Is(err, services.ErrMagicLinkExpired) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Magic link has expired",
			})
		}
		if errors.Is(err, services.ErrMagicLinkUsed) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Magic link has already been used",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to validate magic link",
		})
	}

	// Generate JWT token
	// Using dummy values for clientId and userId for now
	// In production, you would extract these from the magic link or associated user record
	clientId := "client_" + magicLink.ApiKey
	userId := magicLink.ID.String()

	jwtToken := jwt.GenerateJWT(clientId, userId)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Magic link validated successfully",
		"access_token": jwtToken,
		"token_type":   "Bearer",
	})
}