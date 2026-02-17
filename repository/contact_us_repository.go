package repository

import (
	"dhl/models"

	"gorm.io/gorm"
)

type ContactRepository interface {
	Save(response *models.ContactUsResponse) error
	GetAll() ([]models.ContactUsResponse, error)
}

type ContactRepoImpl struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &ContactRepoImpl{db}
}

func (r *ContactRepoImpl) Save(response *models.ContactUsResponse) error {
	return r.db.Create(response).Error
}

func (r *ContactRepoImpl) GetAll() ([]models.ContactUsResponse, error) {
	var responses []models.ContactUsResponse
	err := r.db.Order("created_at DESC").Find(&responses).Error
	return responses, err
}
