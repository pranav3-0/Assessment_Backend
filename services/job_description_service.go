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
	UpdateJob(job models.JobDescription, modifiedBy string) error
DeleteJob(jobID int64, modifiedBy string) error
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

func (s *JobDescriptionServiceImpl) UpdateJob(job models.JobDescription, modifiedBy string) error {

	tx := s.db.Begin()

	existing, err := s.repo.GetByID(job.JobID)
	if err != nil {
		tx.Rollback()
		return err
	}

	existing.Title = job.Title
	existing.Description = job.Description
	existing.RequiredSkills = job.RequiredSkills
	existing.Level = job.Level
	existing.ModifiedOn = time.Now()
	existing.ModifiedBy = modifiedBy

	if err := s.repo.Update(tx, existing); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *JobDescriptionServiceImpl) DeleteJob(jobID int64, modifiedBy string) error {

	tx := s.db.Begin()

	if err := s.repo.SoftDelete(tx, jobID, modifiedBy); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}