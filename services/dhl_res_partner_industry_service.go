package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"time"
)

type DHLResPartnerIndustryService interface {
	Create(ctx context.Context, req models.DHLResPartnerIndustry) error
	List(ctx context.Context) ([]models.DHLResPartnerIndustry, error)
	Update(ctx context.Context, req models.DHLResPartnerIndustry) error
	Delete(ctx context.Context, id int64) error
}

type DHLResPartnerIndustryServiceImpl struct {
	repo repository.DHLResPartnerIndustryRepository
}

func NewDHLResPartnerIndustryService(repo repository.DHLResPartnerIndustryRepository) DHLResPartnerIndustryService {
	return &DHLResPartnerIndustryServiceImpl{repo}
}

func (s *DHLResPartnerIndustryServiceImpl) Create(ctx context.Context, req models.DHLResPartnerIndustry) error {
	return s.repo.Create(ctx, &req)
}

func (s *DHLResPartnerIndustryServiceImpl) List(ctx context.Context) ([]models.DHLResPartnerIndustry, error) {
	return s.repo.List(ctx)
}

func (s *DHLResPartnerIndustryServiceImpl) Update(ctx context.Context, req models.DHLResPartnerIndustry) error {
	existing, err := s.repo.GetByID(ctx, req.PartnerIndustryID)
	if err != nil {
		return err
	}

	existing.Name = req.Name
	existing.FullName = req.FullName
	existing.UpdatedBy = req.UpdatedBy
	existing.UpdatedAt = time.Now()
	return s.repo.Update(ctx, existing)
}

func (s *DHLResPartnerIndustryServiceImpl) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
