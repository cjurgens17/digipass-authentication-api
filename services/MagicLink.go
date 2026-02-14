package services

import (
	"DigiPassAuthenticationApi/packages/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MagicLinkService struct {
	db *gorm.DB
}

type MagicLinkResult struct {
	Token      string
	MagicLink  *models.MagicLink
	MagicLinkURL string
}

func NewMagicLinkService(db *gorm.DB) *MagicLinkService {
	return &MagicLinkService{db: db}
}

// generateSecureToken generates a cryptographically secure random token
func (s *MagicLinkService) generateSecureToken() (string, error) {
	// Generate 32 bytes (256 bits) of random data
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	// Encode to URL-safe base64
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	return token, nil
}

// hashToken creates a SHA-256 hash of the token for storage
func (s *MagicLinkService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// VerifyMagicLink creates a new magic link record with a secure token
func (s *MagicLinkService) VerifyMagicLink(apiKey string, expirationMinutes uint8, callbackBaseURL string) (*MagicLinkResult, error) {
	// Validate apiKey exists (in a real scenario, you'd validate against a stored API key)
	if apiKey == "" {
		return nil, ErrInvalidApiKey
	}

	// Generate secure token
	token, err := s.generateSecureToken()
	if err != nil {
		return nil, err
	}

	// Hash the token for storage
	tokenHash := s.hashToken(token)

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)

	// Create magic link record
	magicLink := &models.MagicLink{
		ApiKey:    apiKey,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	// Save to database
	if err := s.db.Create(magicLink).Error; err != nil {
		return nil, fmt.Errorf("failed to create magic link: %w", err)
	}

	// Construct the magic link URL
	magicLinkURL := fmt.Sprintf("%s?token=%s", callbackBaseURL, token)

	return &MagicLinkResult{
		Token:        token,
		MagicLink:    magicLink,
		MagicLinkURL: magicLinkURL,
	}, nil
}

// ValidateToken validates a magic link token and marks it as used
func (s *MagicLinkService) ValidateToken(token string) (*models.MagicLink, error) {
	// Hash the provided token
	tokenHash := s.hashToken(token)

	// Find the magic link by token hash
	var magicLink models.MagicLink
	err := s.db.Where("token_hash = ?", tokenHash).First(&magicLink).Error

	if err == gorm.ErrRecordNotFound {
		return nil, ErrMagicLinkNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if already used
	if magicLink.UsedAt != nil {
		return nil, ErrMagicLinkUsed
	}

	// Check if expired
	if time.Now().After(magicLink.ExpiresAt) {
		return nil, ErrMagicLinkExpired
	}

	// Mark as used
	now := time.Now()
	magicLink.UsedAt = &now

	if err := s.db.Save(&magicLink).Error; err != nil {
		return nil, fmt.Errorf("failed to mark magic link as used: %w", err)
	}

	return &magicLink, nil
}
