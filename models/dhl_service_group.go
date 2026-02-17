package models

import "time"

type DHLServiceGroup struct {
	ServiceGrpID int64     `gorm:"primaryKey;column:service_grp_id"`
	Name         string    `gorm:"column:name"`
	CreatedBy    *int      `gorm:"column:created_by"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedBy    *int      `gorm:"column:updated_by"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	Active       *bool     `gorm:"column:active"`
}

func (DHLServiceGroup) TableName() string {
	return "dhl_service_group"
}
