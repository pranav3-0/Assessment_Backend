package repository

import (
	"context"
	"dhl/models"
	"dhl/utils"
	"errors"

	"gorm.io/gorm"
)

type DHLResPartnerIndustryRepository interface {
	Create(ctx context.Context, ind *models.DHLResPartnerIndustry) error
	List(ctx context.Context) ([]models.DHLResPartnerIndustry, error)
	GetByID(ctx context.Context, id int64) (*models.DHLResPartnerIndustry, error)
	Update(ctx context.Context, ind *models.DHLResPartnerIndustry) error
	Delete(ctx context.Context, id int64) error
}

type DHLResPartnerIndustryRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLResPartnerIndustryRepository(db *gorm.DB) DHLResPartnerIndustryRepository {
	return &DHLResPartnerIndustryRepositoryImpl{db}
}

func (r *DHLResPartnerIndustryRepositoryImpl) Create(ctx context.Context, ind *models.DHLResPartnerIndustry) error {
	return r.db.WithContext(ctx).Create(ind).Error
}

func (r *DHLResPartnerIndustryRepositoryImpl) List(ctx context.Context) ([]models.DHLResPartnerIndustry, error) {
	var list []models.DHLResPartnerIndustry
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("partner_industry_id DESC").
		Find(&list).Error
	return list, err
}

func (r *DHLResPartnerIndustryRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.DHLResPartnerIndustry, error) {
	var ind models.DHLResPartnerIndustry
	err := r.db.WithContext(ctx).
		Where("partner_industry_id = ? AND active = true", id).
		First(&ind).Error
	if err != nil {
		return nil, err
	}
	return &ind, nil
}

func (r *DHLResPartnerIndustryRepositoryImpl) Update(ctx context.Context, ind *models.DHLResPartnerIndustry) error {
	if ind.PartnerIndustryID == 0 {
		return errors.New("partner_industry_id required for update")
	}

	updateMap := utils.BuildUpdateMap(ind)

	return r.db.Model(&models.DHLResPartnerIndustry{}).
		Where("partner_industry_id = ?", ind.PartnerIndustryID).
		Updates(updateMap).Error
}

func (r *DHLResPartnerIndustryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&models.DHLResPartnerIndustry{}).
		Where("partner_industry_id = ?", id).
		Update("active", false).Error
}
