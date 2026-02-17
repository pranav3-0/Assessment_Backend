package repository

import (
	"context"
	"errors"

	"dhl/models"

	"gorm.io/gorm"
)

type DHLServiceLineRepository interface {
	Create(ctx context.Context, sl *models.DHLServiceLine) error
	GetByID(ctx context.Context, id int64) (*models.DHLServiceLine, error)
	Update(ctx context.Context, sl *models.DHLServiceLine) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]models.DHLServiceLine, error)
}

type DHLServiceLineRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLServiceLineRepository(db *gorm.DB) DHLServiceLineRepository {
	return &DHLServiceLineRepositoryImpl{db: db}
}

func (r *DHLServiceLineRepositoryImpl) Create(ctx context.Context, sl *models.DHLServiceLine) error {
	return r.db.WithContext(ctx).Create(sl).Error
}

func (r *DHLServiceLineRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.DHLServiceLine, error) {
	var sl models.DHLServiceLine
	err := r.db.WithContext(ctx).First(&sl, "service_line_id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &sl, nil
}

func (r *DHLServiceLineRepositoryImpl) Update(ctx context.Context, sl *models.DHLServiceLine) error {
	return r.db.WithContext(ctx).Save(sl).Error
}

func (r *DHLServiceLineRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&models.DHLServiceLine{}).
		Where("service_line_id = ?", id).
		Update("active", false).Error
}

func (r *DHLServiceLineRepositoryImpl) List(ctx context.Context) ([]models.DHLServiceLine, error) {
	var list []models.DHLServiceLine
	err := r.db.WithContext(ctx).
		Where("active = true").
		Order("service_line_id desc").
		Find(&list).Error
	return list, err
}
