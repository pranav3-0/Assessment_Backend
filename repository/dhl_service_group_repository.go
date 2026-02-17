package repository

import (
	"context"
	"dhl/models"
	"dhl/utils"
	"errors"

	"gorm.io/gorm"
)

type DHLServiceGroupRepository interface {
	Create(ctx context.Context, grp *models.DHLServiceGroup) error
	List(ctx context.Context) ([]models.DHLServiceGroup, error)
	GetByID(ctx context.Context, id int64) (*models.DHLServiceGroup, error)
	Update(ctx context.Context, grp *models.DHLServiceGroup) error
	Delete(ctx context.Context, id int64) error
}

type DHLServiceGroupRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLServiceGroupRepository(db *gorm.DB) DHLServiceGroupRepository {
	return &DHLServiceGroupRepositoryImpl{db}
}

func (r *DHLServiceGroupRepositoryImpl) Create(ctx context.Context, grp *models.DHLServiceGroup) error {
	return r.db.WithContext(ctx).Create(grp).Error
}

func (r *DHLServiceGroupRepositoryImpl) List(ctx context.Context) ([]models.DHLServiceGroup, error) {
	var list []models.DHLServiceGroup
	err := r.db.WithContext(ctx).
		Where("active = true").
		Order("service_grp_id DESC").
		Find(&list).Error
	return list, err
}

func (r *DHLServiceGroupRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.DHLServiceGroup, error) {
	var grp models.DHLServiceGroup
	err := r.db.WithContext(ctx).
		Where("service_grp_id = ? AND active = true", id).
		First(&grp).Error

	if err != nil {
		return nil, err
	}
	return &grp, nil
}

func (r *DHLServiceGroupRepositoryImpl) Update(ctx context.Context, grp *models.DHLServiceGroup) error {
	if grp.ServiceGrpID == 0 {
		return errors.New("service_grp_id required for update")
	}

	updateMap := utils.BuildUpdateMap(grp)

	return r.db.Model(&models.DHLServiceGroup{}).
		Where("service_grp_id = ?", grp.ServiceGrpID).
		Updates(updateMap).Error
}

func (r *DHLServiceGroupRepositoryImpl) Delete(ctx context.Context, id int64) error {
	active := false
	return r.db.WithContext(ctx).
		Model(&models.DHLServiceGroup{}).
		Where("service_grp_id = ?", id).
		Update("active", &active).Error
}
