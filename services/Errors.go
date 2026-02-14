package services

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrRecordNotFound = errors.New("record not found")
	ErrMagicLinkExpired = errors.New("magic link has expired")
	ErrMagicLinkUsed = errors.New("magic link has already been used")
	ErrMagicLinkNotFound = errors.New("magic link not found")
	ErrInvalidApiKey = errors.New("invalid api key")
)