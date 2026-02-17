package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"time"
)

type DHLSubBusinessPartnerService interface {
	Create(ctx context.Context, req models.DHLSubBusinessPartner) error
	List(ctx context.Context) ([]models.DHLSubBusinessPartner, error)
	Update(ctx context.Context, req models.DHLSubBusinessPartner) error
	Delete(ctx context.Context, id int64) error
}

type DHLSubBusinessPartnerServiceImpl struct {
	repo repository.DHLSubBusinessPartnerRepository
}

func NewDHLSubBusinessPartnerService(repo repository.DHLSubBusinessPartnerRepository) DHLSubBusinessPartnerService {
	return &DHLSubBusinessPartnerServiceImpl{repo}
}

func (s *DHLSubBusinessPartnerServiceImpl) Create(ctx context.Context, req models.DHLSubBusinessPartner) error {
	return s.repo.Create(ctx, &req)
}

func (s *DHLSubBusinessPartnerServiceImpl) List(ctx context.Context) ([]models.DHLSubBusinessPartner, error) {
	return s.repo.List(ctx)
}

func (s *DHLSubBusinessPartnerServiceImpl) Update(ctx context.Context, req models.DHLSubBusinessPartner) error {
	existing, err := s.repo.GetByID(ctx, req.SubBusinessPartnerID)
	if err != nil {
		return err
	}
	existing.Name = req.Name
	existing.XmlID = req.XmlID
	existing.LineID = req.LineID
	existing.UpdatedBy = req.UpdatedBy
	existing.UpdatedAt = time.Now()

	return s.repo.Update(ctx, existing)
}

func (s *DHLSubBusinessPartnerServiceImpl) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
