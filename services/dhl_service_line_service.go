package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"time"
)

type DHLServiceLineService interface {
	Create(ctx context.Context, sl models.DHLServiceLine) error
	Update(ctx context.Context, sl models.DHLServiceLine) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]models.DHLServiceLine, error)
	GetByID(ctx context.Context, id int64) (*models.DHLServiceLine, error)
}

type DHLServiceLineServiceImpl struct {
	repo repository.DHLServiceLineRepository
}

func NewDHLServiceLineService(repo repository.DHLServiceLineRepository) DHLServiceLineService {
	return &DHLServiceLineServiceImpl{repo: repo}
}

func (s *DHLServiceLineServiceImpl) Create(ctx context.Context, sl models.DHLServiceLine) error {
	active := true
	sl.CreatedAt = time.Now()
	sl.Active = &active
	return s.repo.Create(ctx, &sl)
}

func (s *DHLServiceLineServiceImpl) Update(ctx context.Context, sl models.DHLServiceLine) error {
	existing, err := s.repo.GetByID(ctx, sl.ServiceLineID)
	if err != nil {
		return err
	}

	existing.Name = sl.Name
	existing.XmlID = sl.XmlID
	existing.UpdatedBy = sl.UpdatedBy
	existing.UpdatedAt = time.Now()
	return s.repo.Update(ctx, existing)
}

func (s *DHLServiceLineServiceImpl) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *DHLServiceLineServiceImpl) List(ctx context.Context) ([]models.DHLServiceLine, error) {
	return s.repo.List(ctx)
}

func (s *DHLServiceLineServiceImpl) GetByID(ctx context.Context, id int64) (*models.DHLServiceLine, error) {
	return s.repo.GetByID(ctx, id)
}
