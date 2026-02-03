package services

import (
	"DigiPassAuthenticationApi/packages/models"
	"errors"
	"gorm.io/gorm"
)

type TenantService struct {
	db *gorm.DB
}

func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{
		db: db,
	}
}

func (s *TenantService) CreateUniqueTenantSlug() (string,error){
	maxAttempts := 10

	for i := 0; i < maxAttempts; i++ {
		slug := models.Tenant{}.CreateSlug()

		var count int64
		err := s.db.Model(&models.Tenant{}).Where("slug = ?", slug).Count(&count).Error
		if err != nil {
			return "", err
		}

		if count == 0 {
			return slug, nil
		}
	}
	return "", errors.New("Failed to generate unique slug after multiple attempts")
}