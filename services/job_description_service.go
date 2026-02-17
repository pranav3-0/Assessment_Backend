package services

import (
	"dhl/models"
	"dhl/repository"
	"time"

	"gorm.io/gorm"
)

type JobDescriptionService interface {
	CreateJob(job models.JobDescription, createdBy string) error
	GetJobs(limit, offset int) ([]models.JobDescription, int64, error)
}

type JobDescriptionServiceImpl struct {
	repo repository.JobDescriptionRepository
	db   *gorm.DB
}

func NewJobDescriptionService(repo repository.JobDescriptionRepository, db *gorm.DB) JobDescriptionService {
	return &JobDescriptionServiceImpl{repo: repo, db: db}
}

func (s *JobDescriptionServiceImpl) CreateJob(job models.JobDescription, createdBy string) error {

	job.CreatedOn = time.Now()
	job.ModifiedOn = time.Now()
	job.CreatedBy = createdBy
	job.ModifiedBy = createdBy
	job.IsActive = true
	job.IsDeleted = false

	tx := s.db.Begin()

	if err := s.repo.Create(tx, &job); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *JobDescriptionServiceImpl) GetJobs(limit, offset int) ([]models.JobDescription, int64, error) {
	return s.repo.GetAll(limit, offset)
}
