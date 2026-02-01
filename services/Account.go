package services

import (
	"DigiPassAuthenticationApi/packages/models"
	"errors"
	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

func (s *AccountService) CreateAccount(name string, email string) (*models.Account, error) {
	if name == "" || email == "" {
		return nil, errors.New("Name and email are required")
	}

	//Open a transaction for rollbacks
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//1. Create Account
	account := &models.Account{
		Name:   name,
		Email:  email,
		Status: "active",
	}

	if err := tx.Create(account).Error; err != nil {
		return nil, err
	}

	//Tenant Service
	ts := NewTenantService(s.db)

	//2. Create Tenant
	slug, err := ts.CreateUniqueTenantSlug(s.db)
	if err != nil {
		return nil, err
	}

	tenant := &models.Tenant {
		AccountID: account.ID,
		Slug: slug, 
	}

	if err := tx.Create(tenant).Error; err != nil {
		return nil, err
	}

	//3. Link Creating User to AccountUsers
	//4. Commit transaction
	return account, nil
}
