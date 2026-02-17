package models

import "time"

type Client struct {
	ClientID   string    `gorm:"column:client_id;primaryKey"`
	ClientName string    `gorm:"column:client_name"`
	AuthType   string    `gorm:"column:auth_type"`
	AuthConfig string    `gorm:"column:auth_config"`
	IsActive   bool      `gorm:"column:is_active"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (Client) TableName() string {
	return "client_mst"
}
