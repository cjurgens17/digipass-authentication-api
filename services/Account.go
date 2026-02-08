package services

import (
	"DigiPassAuthenticationApi/packages/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) CreateAccount(name string, email string) (*models.Account, error) {
	if name == "" || email == "" {
		return nil, errors.New("Name and email are required")
	}

	//1. Create Account
	account := &models.Account{
		Name:   name,
		Email:  email,
		Status: "active",
	}

	if err := s.db.Create(account).Error; err != nil {
		return nil, err
	}

	//Tenant Service
	ts := NewTenantService(s.db)

	//2. Create Tenant
	slug, err := ts.CreateUniqueTenantSlug()
	if err != nil {
		return nil, err
	}

	tenant := &models.Tenant{
		AccountID: account.ID,
		Slug:      slug,
		Name:      fmt.Sprintf("%s-%s", account.Name, slug),
	}

	if err := s.db.Create(tenant).Error; err != nil {
		return nil, err
	}

	//3. Link User to AccountUser Table
	aus := NewAccountUsersService(s.db)
	accountUserErr := aus.CreateUser(account.ID, account.Email, "abc123", "owner")
	if accountUserErr != nil {
		return nil, accountUserErr
	}

	return account, nil
}
