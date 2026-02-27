package repository

import (
	"dhl/models"
	"gorm.io/gorm"
	"time"
	
)

type JobDescriptionRepository interface {
	Create(tx *gorm.DB, job *models.JobDescription) error
	GetAll(limit, offset int) ([]models.JobDescription, int64, error)
	Update(tx *gorm.DB, job *models.JobDescription) error
	SoftDelete(tx *gorm.DB, jobID int64, modifiedBy string) error
	GetByID(jobID int64) (*models.JobDescription, error)
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

func (r *JobDescriptionRepositoryImpl) GetByID(jobID int64) (*models.JobDescription, error) {
	var job models.JobDescription
	err := r.db.
		Where("job_id = ? AND is_deleted = false", jobID).
		First(&job).Error

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *JobDescriptionRepositoryImpl) Update(tx *gorm.DB, job *models.JobDescription) error {
	return tx.Save(job).Error
}

func (r *JobDescriptionRepositoryImpl) SoftDelete(tx *gorm.DB, jobID int64, modifiedBy string) error {
	return tx.Model(&models.JobDescription{}).
		Where("job_id = ? AND is_deleted = false", jobID).
		Updates(map[string]interface{}{
			"is_deleted":  true,
			"is_active":   false,
			"modified_on": time.Now(),
			"modified_by": modifiedBy,
		}).Error
}