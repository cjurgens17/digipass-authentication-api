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

	//Open a transaction
	tx := s.db.Begin()

	//1. Create Account
	account := &models.Account{
		Name:   name,
		Email:  email,
		Status: "active",
	}

	if err := tx.Create(account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//Tenant Service
	ts := NewTenantService(s.db)

	//2. Create Tenant
	slug, err := ts.CreateUniqueTenantSlug()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tenant := &models.Tenant{
		AccountID: account.ID,
		Slug:      slug,
	}

	if err := tx.Create(tenant).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	//3. Link User to AccountUser Table
	aus := NewAccountUsersService(s.db)
	accountUserErr := aus.CreateUser(account.ID, account.Email, "abc123", "owner")
	if accountUserErr != nil {
		tx.Rollback()
		return nil, accountUserErr
	}

	//4. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return account, nil
}
