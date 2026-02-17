package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"time"
)

type DHLServiceService interface {
	Create(ctx context.Context, req models.DHLService) error
	List(ctx context.Context) ([]models.DHLService, error)
	Update(ctx context.Context, req models.DHLService) error
	Delete(ctx context.Context, id int64) error
}

type DHLServiceServiceImpl struct {
	repo repository.DHLServiceRepository
}

func NewDHLServiceService(repo repository.DHLServiceRepository) DHLServiceService {
	return &DHLServiceServiceImpl{repo}
}

func (s *DHLServiceServiceImpl) Create(ctx context.Context, req models.DHLService) error {
	return s.repo.Create(ctx, &req)
}

func (s *DHLServiceServiceImpl) List(ctx context.Context) ([]models.DHLService, error) {
	return s.repo.List(ctx)
}

func (s *DHLServiceServiceImpl) Update(ctx context.Context, req models.DHLService) error {
	existing, err := s.repo.GetByID(ctx, req.ServiceID)
	if err != nil {
		return err
	}

	existing.ServiceName = req.ServiceName
	existing.UpdatedBy = req.UpdatedBy
	existing.UpdatedAt = time.Now()

	return s.repo.Update(ctx, existing)
}

func (s *DHLServiceServiceImpl) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
