package models

import "time"

type DHLBusinessPartner struct {
	BusinessPartnerID int64     `json:"business_partner_id" gorm:"primaryKey;column:business_partner_id"`
	Name              string    `json:"name" gorm:"column:name"`
	XMLID             string    `json:"xml_id" gorm:"column:xml_id"`
	CreatedBy         int       `json:"created_by" gorm:"column:created_by"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedBy         int       `json:"updated_by" gorm:"column:updated_by"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at"`
	Active            bool      `json:"active" gorm:"column:active"`
}

func (DHLBusinessPartner) TableName() string {
	return "dhl_business_partner"
}
