package models

import "time"

type DHLServiceLine struct {
	ServiceLineID int64     `gorm:"primaryKey;column:service_line_id" json:"service_line_id"`
	Name          string    `gorm:"column:name" json:"name"`
	XmlID         *int      `gorm:"column:xml_id" json:"xml_id"`
	CreatedBy     *int      `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int      `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
	Active        *bool     `gorm:"column:active" json:"active"`
}

func (DHLServiceLine) TableName() string {
	return "dhl_service_line"
}
