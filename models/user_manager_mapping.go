package models

import "time"

type UserManagerMapping struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	ManagerID string    `gorm:"column:manager_id" json:"manager_id"`
	IsActive  bool      `gorm:"column:is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy string    `gorm:"column:updated_by" json:"updated_by"`
}

func (UserManagerMapping) TableName() string {
	return "user_manager_mapping"
}
