package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"time"
)

type DHLCenterService interface {
	CreateCenter(ctx context.Context, center models.DHLCenter) error
	ListCenters(ctx context.Context) ([]models.DHLCenter, error)
	UpdateCenter(ctx context.Context, center models.DHLCenter) error
	DeleteCenter(ctx context.Context, id int64) error
}

type DHLCenterServiceImpl struct {
	repo repository.DHLCenterRepository
}

func NewDHLCenterService(repo repository.DHLCenterRepository) DHLCenterService {
	return &DHLCenterServiceImpl{repo}
}

func (s *DHLCenterServiceImpl) CreateCenter(ctx context.Context, center models.DHLCenter) error {
	return s.repo.Create(ctx, &center)
}

func (s *DHLCenterServiceImpl) ListCenters(ctx context.Context) ([]models.DHLCenter, error) {
	return s.repo.List(ctx)
}

func (s *DHLCenterServiceImpl) UpdateCenter(ctx context.Context, center models.DHLCenter) error {
	updates := &models.DHLCenter{
		CenterID:   center.CenterID,
		CenterName: center.CenterName,
		UpdatedAt:  time.Now(),
	}
	return s.repo.Update(ctx, updates)
}

func (s *DHLCenterServiceImpl) DeleteCenter(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
