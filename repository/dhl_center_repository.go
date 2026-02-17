package repository

import (
	"context"
	"dhl/models"
	"dhl/utils"
	"errors"

	"gorm.io/gorm"
)

type DHLCenterRepository interface {
	Create(ctx context.Context, center *models.DHLCenter) error
	List(ctx context.Context) ([]models.DHLCenter, error)
	Update(ctx context.Context, updates *models.DHLCenter) error
	Delete(ctx context.Context, id int64) error
}

type DHLCenterRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLCenterRepository(db *gorm.DB) DHLCenterRepository {
	return &DHLCenterRepositoryImpl{db}
}

func (r *DHLCenterRepositoryImpl) Create(ctx context.Context, center *models.DHLCenter) error {
	return r.db.WithContext(ctx).Create(center).Error
}

func (r *DHLCenterRepositoryImpl) List(ctx context.Context) ([]models.DHLCenter, error) {
	var centers []models.DHLCenter
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("center_id DESC").
		Find(&centers).Error
	return centers, err
}

func (r *DHLCenterRepositoryImpl) Update(ctx context.Context, updates *models.DHLCenter) error {
	if updates.CenterID == 0 {
		return errors.New("assessment_sequence required for update")
	}
	updateMap := utils.BuildUpdateMap(updates)
	return r.db.Model(&models.DHLCenter{}).
		Where("center_id = ?", updates.CenterID).
		Updates(updateMap).Error
}

func (r *DHLCenterRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&models.DHLCenter{}).
		Where("center_id = ?", id).
		Update("active", false).Error
}
