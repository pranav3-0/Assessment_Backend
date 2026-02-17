package models

import "time"

type DHLService struct {
	ServiceID   int64     `gorm:"primaryKey;column:service_id"`
	ServiceName string    `gorm:"column:service_name"`
	CreatedBy   *int      `gorm:"column:created_by"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedBy   *int      `gorm:"column:updated_by"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	Active      *bool     `gorm:"column:active"`
}

func (DHLService) TableName() string {
	return "dhl_service"
}
