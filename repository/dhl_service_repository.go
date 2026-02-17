package repository

import (
	"context"
	"dhl/models"
	"dhl/utils"
	"errors"

	"gorm.io/gorm"
)

type DHLServiceRepository interface {
	Create(ctx context.Context, service *models.DHLService) error
	List(ctx context.Context) ([]models.DHLService, error)
	GetByID(ctx context.Context, id int64) (*models.DHLService, error)
	Update(ctx context.Context, service *models.DHLService) error
	Delete(ctx context.Context, id int64) error
}

type DHLServiceRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLServiceRepository(db *gorm.DB) DHLServiceRepository {
	return &DHLServiceRepositoryImpl{db}
}

func (r *DHLServiceRepositoryImpl) Create(ctx context.Context, service *models.DHLService) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *DHLServiceRepositoryImpl) List(ctx context.Context) ([]models.DHLService, error) {
	var list []models.DHLService
	err := r.db.WithContext(ctx).
		Where("active = true").
		Order("service_id DESC").
		Find(&list).Error
	return list, err
}

func (r *DHLServiceRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.DHLService, error) {
	var service models.DHLService
	err := r.db.WithContext(ctx).
		Where("service_id = ? AND active = true", id).
		First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *DHLServiceRepositoryImpl) Update(ctx context.Context, service *models.DHLService) error {
	if service.ServiceID == 0 {
		return errors.New("service_id required for update")
	}

	updateMap := utils.BuildUpdateMap(service)

	return r.db.Model(&models.DHLService{}).
		Where("service_id = ?", service.ServiceID).
		Updates(updateMap).Error
}

func (r *DHLServiceRepositoryImpl) Delete(ctx context.Context, id int64) error {
	active := false
	return r.db.WithContext(ctx).
		Model(&models.DHLService{}).
		Where("service_id = ?", id).
		Update("active", &active).Error
}
