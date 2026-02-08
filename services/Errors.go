package services

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrRecordNotFound = errors.New("record not found")
)