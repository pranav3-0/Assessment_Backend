package repository

import (
	"context"
	"dhl/models"
	"dhl/utils"
	"errors"

	"gorm.io/gorm"
)

type DHLBusinessPartnerRepository interface {
	Create(ctx context.Context, bp *models.DHLBusinessPartner) error
	List(ctx context.Context) ([]models.DHLBusinessPartner, error)
	Update(ctx context.Context, updates *models.DHLBusinessPartner) error
	Delete(ctx context.Context, id int64) error
}

type DHLBusinessPartnerRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLBusinessPartnerRepository(db *gorm.DB) DHLBusinessPartnerRepository {
	return &DHLBusinessPartnerRepositoryImpl{db}
}

func (r *DHLBusinessPartnerRepositoryImpl) Create(ctx context.Context, bp *models.DHLBusinessPartner) error {
	return r.db.WithContext(ctx).Create(bp).Error
}

func (r *DHLBusinessPartnerRepositoryImpl) List(ctx context.Context) ([]models.DHLBusinessPartner, error) {
	var partners []models.DHLBusinessPartner
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("business_partner_id DESC").
		Find(&partners).Error
	return partners, err
}

func (r *DHLBusinessPartnerRepositoryImpl) Update(ctx context.Context, updates *models.DHLBusinessPartner) error {
	if updates.BusinessPartnerID == 0 {
		return errors.New("business_partner_id required for update")
	}
	updateMap := utils.BuildUpdateMap(updates)
	return r.db.Model(&models.DHLBusinessPartner{}).
		Where("business_partner_id = ?", updates.BusinessPartnerID).
		Updates(updateMap).Error
}

func (r *DHLBusinessPartnerRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&models.DHLBusinessPartner{}).
		Where("business_partner_id = ?", id).
			Update("active", false).Error
}
