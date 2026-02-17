package repository

import (
	"context"
	"dhl/models"
	"dhl/utils"
	"errors"

	"gorm.io/gorm"
)

type DHLSubServiceRepository interface {
	Create(ctx context.Context, ss *models.DHLSubService) error
	List(ctx context.Context) ([]models.DHLSubService, error)
	GetByID(ctx context.Context, id int64) (*models.DHLSubService, error)
	Update(ctx context.Context, ss *models.DHLSubService) error
	Delete(ctx context.Context, id int64) error
}

type DHLSubServiceRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLSubServiceRepository(db *gorm.DB) DHLSubServiceRepository {
	return &DHLSubServiceRepositoryImpl{db}
}

func (r *DHLSubServiceRepositoryImpl) Create(ctx context.Context, ss *models.DHLSubService) error {
	return r.db.WithContext(ctx).Create(ss).Error
}

func (r *DHLSubServiceRepositoryImpl) List(ctx context.Context) ([]models.DHLSubService, error) {
	var list []models.DHLSubService
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("sub_service_id DESC").
		Find(&list).Error
	return list, err
}

func (r *DHLSubServiceRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.DHLSubService, error) {
	var ss models.DHLSubService

	err := r.db.WithContext(ctx).
		Where("sub_service_id = ? AND active = true", id).
		First(&ss).Error

	if err != nil {
		return nil, err
	}

	return &ss, nil
}

func (r *DHLSubServiceRepositoryImpl) Update(ctx context.Context, ss *models.DHLSubService) error {
	if ss.SubServiceID == 0 {
		return errors.New("sub_service_id required for update")
	}

	updateMap := utils.BuildUpdateMap(ss)

	return r.db.Model(&models.DHLSubService{}).
		Where("sub_service_id = ?", ss.SubServiceID).
		Updates(updateMap).Error
}

func (r *DHLSubServiceRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&models.DHLSubService{}).
		Where("sub_service_id = ?", id).
		Update("active", false).Error
}
