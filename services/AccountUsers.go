package services

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"DigiPassAuthenticationApi/packages/models"
	"errors"
	"fmt"
)

type AccountUsersService struct {
	db *gorm.DB
}

func NewAccountUsersService(db *gorm.DB) *AccountUsersService{
	return &AccountUsersService{db:db}
}

func isValidRole(role string) bool {
    return role == "owner" || role == "admin" || role == "member"
}

func (s *AccountUsersService) CreateUser(accountID uuid.UUID, email string, passwordHash string, role string) (error) {
	if email == "" || !isValidRole(role) {
		return errors.New("Required Inputs do not match for AccountUsersService.CreateUser")
	}

	newAccountUser := &models.AccountUser{
		AccountID: accountID,
		Email: email,
		PasswordHash: passwordHash,
		Role: role,
	}


	if err := s.db.Create(newAccountUser).Error; err != nil {
		return fmt.Errorf("failed to create account user: %w", err)
	}

	return nil
}