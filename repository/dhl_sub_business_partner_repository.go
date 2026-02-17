package repository

import (
	"context"
	"dhl/models"
	"dhl/utils"
	"errors"

	"gorm.io/gorm"
)

type DHLSubBusinessPartnerRepository interface {
	Create(ctx context.Context, sbp *models.DHLSubBusinessPartner) error
	List(ctx context.Context) ([]models.DHLSubBusinessPartner, error)
	GetByID(ctx context.Context, id int64) (*models.DHLSubBusinessPartner, error)
	Update(ctx context.Context, sbp *models.DHLSubBusinessPartner) error
	Delete(ctx context.Context, id int64) error
}

type DHLSubBusinessPartnerRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLSubBusinessPartnerRepository(db *gorm.DB) DHLSubBusinessPartnerRepository {
	return &DHLSubBusinessPartnerRepositoryImpl{db}
}

func (r *DHLSubBusinessPartnerRepositoryImpl) Create(ctx context.Context, sbp *models.DHLSubBusinessPartner) error {
	return r.db.WithContext(ctx).Create(sbp).Error
}

func (r *DHLSubBusinessPartnerRepositoryImpl) List(ctx context.Context) ([]models.DHLSubBusinessPartner, error) {
	var list []models.DHLSubBusinessPartner
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("sub_business_partner_id DESC").
		Find(&list).Error
	return list, err
}

func (r *DHLSubBusinessPartnerRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.DHLSubBusinessPartner, error) {
	var sbp models.DHLSubBusinessPartner
	err := r.db.WithContext(ctx).
		Where("sub_business_partner_id = ? AND active = true", id).
		First(&sbp).Error

	if err != nil {
		return nil, err
	}

	return &sbp, nil
}

func (r *DHLSubBusinessPartnerRepositoryImpl) Update(ctx context.Context, sbp *models.DHLSubBusinessPartner) error {
	if sbp.SubBusinessPartnerID == 0 {
		return errors.New("sub_business_partner_id required for update")
	}

	updateMap := utils.BuildUpdateMap(sbp)

	return r.db.Model(&models.DHLSubBusinessPartner{}).
		Where("sub_business_partner_id = ?", sbp.SubBusinessPartnerID).
		Updates(updateMap).Error
}

func (r *DHLSubBusinessPartnerRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&models.DHLSubBusinessPartner{}).
		Where("sub_business_partner_id = ?", id).
		Update("active", false).Error
}
