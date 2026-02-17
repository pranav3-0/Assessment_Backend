package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"time"
)

type DHLSubServiceService interface {
	Create(ctx context.Context, req models.DHLSubService) error
	List(ctx context.Context) ([]models.DHLSubService, error)
	Update(ctx context.Context, req models.DHLSubService) error
	Delete(ctx context.Context, id int64) error
}

type DHLSubServiceServiceImpl struct {
	repo repository.DHLSubServiceRepository
}

func NewDHLSubServiceService(repo repository.DHLSubServiceRepository) DHLSubServiceService {
	return &DHLSubServiceServiceImpl{repo}
}

func (s *DHLSubServiceServiceImpl) Create(ctx context.Context, req models.DHLSubService) error {
	return s.repo.Create(ctx, &req)
}

func (s *DHLSubServiceServiceImpl) List(ctx context.Context) ([]models.DHLSubService, error) {
	return s.repo.List(ctx)
}

func (s *DHLSubServiceServiceImpl) Update(ctx context.Context, req models.DHLSubService) error {
	existing, err := s.repo.GetByID(ctx, req.SubServiceID)
	if err != nil {
		return err
	}
	existing.Name = req.Name
	existing.UpdatedBy = req.UpdatedBy
	now := time.Now()
	existing.UpdatedAt = &now

	return s.repo.Update(ctx, existing)
}

func (s *DHLSubServiceServiceImpl) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
