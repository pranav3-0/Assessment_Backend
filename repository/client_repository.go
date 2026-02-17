package repository

import (
	"context"
	"dhl/models"

	"gorm.io/gorm"
)

type ClientRepository interface {
	GetByName(ctx context.Context, name string) (*models.Client, error)
}

type ClientRepositoryImpl struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &ClientRepositoryImpl{db: db}
}

func (r *ClientRepositoryImpl) GetByName(ctx context.Context, name string) (*models.Client, error) {
	var client models.Client
	if err := r.db.WithContext(ctx).Where("client_name = ?", name).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}
