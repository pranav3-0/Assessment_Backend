package models

import "time"

type DHLSubBusinessPartner struct {
	SubBusinessPartnerID int64     `gorm:"primaryKey;column:sub_business_partner_id" json:"sub_business_partner_id"`
	Name                 string    `gorm:"column:name" json:"name"`
	LineID               *int      `gorm:"column:line_id" json:"line_id"`
	XmlID                *string   `gorm:"column:xml_id" json:"xml_id"`
	CreatedBy            *int      `gorm:"column:created_by" json:"created_by"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy            *int      `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt            time.Time `gorm:"column:updated_at" json:"updated_at"`
	Active               *bool     `gorm:"column:active" json:"active"`
}

func (DHLSubBusinessPartner) TableName() string {
	return "dhl_sub_business_partner"
}
