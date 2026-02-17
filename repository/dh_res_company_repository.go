package repository

import (
	"context"
	"dhl/models"
	"dhl/utils"
	"errors"

	"gorm.io/gorm"
)

type DHLResCompanyRepository interface {
	Create(ctx context.Context, data *models.DHLResCompany) error
	List(ctx context.Context) ([]models.DHLResCompany, error)
	Update(ctx context.Context, data *models.DHLResCompany) error
	Delete(ctx context.Context, id int64) error
}

type DHLResCompanyRepositoryImpl struct {
	db *gorm.DB
}

func NewDHLResCompanyRepository(db *gorm.DB) DHLResCompanyRepository {
	return &DHLResCompanyRepositoryImpl{db}
}

func (r *DHLResCompanyRepositoryImpl) Create(ctx context.Context, data *models.DHLResCompany) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *DHLResCompanyRepositoryImpl) List(ctx context.Context) ([]models.DHLResCompany, error) {
	var companies []models.DHLResCompany
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("company_id DESC").
		Find(&companies).Error
	return companies, err
}

func (r *DHLResCompanyRepositoryImpl) Update(ctx context.Context, data *models.DHLResCompany) error {
	if data.CompanyID == 0 {
		return errors.New("company_id required for update")
	}
	updates := utils.BuildUpdateMap(data)
	return r.db.WithContext(ctx).Model(models.DHLResCompany{}).
		Where("company_id = ?", data.CompanyID).
		Updates(updates).Error
}

func (r *DHLResCompanyRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&models.DHLResCompany{}).
		Where("company_id = ?", id).
		Update("is_active", false).Error
}
