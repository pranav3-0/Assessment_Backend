package repository

import (
	"dhl/models"
	"gorm.io/gorm"
)

type JobDescriptionRepository interface {
	Create(tx *gorm.DB, job *models.JobDescription) error
	GetAll(limit, offset int) ([]models.JobDescription, int64, error)
}

type JobDescriptionRepositoryImpl struct {
	db *gorm.DB
}

func NewJobDescriptionRepository(db *gorm.DB) JobDescriptionRepository {
	return &JobDescriptionRepositoryImpl{db: db}
}

func (r *JobDescriptionRepositoryImpl) Create(tx *gorm.DB, job *models.JobDescription) error {
	return tx.Create(job).Error
}

func (r *JobDescriptionRepositoryImpl) GetAll(limit, offset int) ([]models.JobDescription, int64, error) {
	var jobs []models.JobDescription
	var total int64

	if err := r.db.Model(&models.JobDescription{}).
		Where("is_deleted = false").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Where("is_deleted = false").
		Order("created_on DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}
