package models

import "time"

type DHLCenter struct {
	CenterID   int64     `json:"center_id" gorm:"primaryKey;column:center_id"`
	CenterName string    `json:"center_name" gorm:"column:center_name"`
	CreatedBy  int       `json:"created_by" gorm:"column:created_by"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedBy  int       `json:"updated_by" gorm:"column:updated_by"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
	Active     bool      `json:"active" gorm:"column:active"`
}

func (DHLCenter) TableName() string {
	return "dhl_center"
}

