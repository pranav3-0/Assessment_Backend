package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"time"
)

type DHLServiceGroupService interface {
	Create(ctx context.Context, req models.DHLServiceGroup) error
	List(ctx context.Context) ([]models.DHLServiceGroup, error)
	Update(ctx context.Context, req models.DHLServiceGroup) error
	Delete(ctx context.Context, id int64) error
}

type DHLServiceGroupServiceImpl struct {
	repo repository.DHLServiceGroupRepository
}

func NewDHLServiceGroupService(repo repository.DHLServiceGroupRepository) DHLServiceGroupService {
	return &DHLServiceGroupServiceImpl{repo}
}

func (s *DHLServiceGroupServiceImpl) Create(ctx context.Context, req models.DHLServiceGroup) error {
	return s.repo.Create(ctx, &req)
}

func (s *DHLServiceGroupServiceImpl) List(ctx context.Context) ([]models.DHLServiceGroup, error) {
	return s.repo.List(ctx)
}

func (s *DHLServiceGroupServiceImpl) Update(ctx context.Context, req models.DHLServiceGroup) error {
	existing, err := s.repo.GetByID(ctx, req.ServiceGrpID)
	if err != nil {
		return err
	}

	existing.Name = req.Name
	existing.UpdatedBy = req.UpdatedBy
	existing.UpdatedAt = time.Now()

	return s.repo.Update(ctx, existing)
}

func (s *DHLServiceGroupServiceImpl) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
