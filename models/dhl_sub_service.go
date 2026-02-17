package models

import "time"

type DHLSubService struct {
	SubServiceID int64      `gorm:"column:sub_service_id;primaryKey;autoIncrement"`
	Name         string     `gorm:"column:name"`
	CreatedBy    *int       `gorm:"column:created_by"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedBy    *int       `gorm:"column:updated_by"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	Active       *bool      `gorm:"column:active"`
}

func (DHLSubService) TableName() string {
	return "dhl_sub_service"
}
