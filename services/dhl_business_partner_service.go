package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
	"time"
)

type DHLBusinessPartnerService interface {
	CreatePartner(ctx context.Context, bp models.DHLBusinessPartner) error
	ListPartners(ctx context.Context) ([]models.DHLBusinessPartner, error)
	UpdatePartner(ctx context.Context, bp models.DHLBusinessPartner) error
	DeletePartner(ctx context.Context, id int64) error
}

type DHLBusinessPartnerServiceImpl struct {
	repo repository.DHLBusinessPartnerRepository
}

func NewDHLBusinessPartnerService(repo repository.DHLBusinessPartnerRepository) DHLBusinessPartnerService {
	return &DHLBusinessPartnerServiceImpl{repo}
}

func (s *DHLBusinessPartnerServiceImpl) CreatePartner(ctx context.Context, bp models.DHLBusinessPartner) error {
	return s.repo.Create(ctx, &bp)
}

func (s *DHLBusinessPartnerServiceImpl) ListPartners(ctx context.Context) ([]models.DHLBusinessPartner, error) {
	return s.repo.List(ctx)
}

func (s *DHLBusinessPartnerServiceImpl) UpdatePartner(ctx context.Context, bp models.DHLBusinessPartner) error {
	updates := &models.DHLBusinessPartner{
		BusinessPartnerID: bp.BusinessPartnerID,
		Name:              bp.Name,
		XMLID:             bp.XMLID,
		UpdatedAt:         time.Now(),
	}
	return s.repo.Update(ctx, updates)
}

func (s *DHLBusinessPartnerServiceImpl) DeletePartner(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
