package services

import (
	"gorm.io/gorm"
	"DigiPassAuthenticationApi/packages/models"
	"errors"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (s *AccountService) CreateAccount(name string, email string) (*models.Account, error){
	if name == "" || email == ""{
		return nil, errors.New("Name and email are required")
	}

	account := &models.Account {
		Name: name,
		Email: email,
		Status: "active",
	}

	if err := s.db.Create(account).Error; err != nil {
		return nil,err
	}
	return account, nil
}