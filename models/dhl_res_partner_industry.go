package models

import "time"

type DHLResPartnerIndustry struct {
	PartnerIndustryID int64     `json:"partner_industry_id" gorm:"primaryKey;column:partner_industry_id"`
	Name              string    `json:"name" gorm:"column:name"`
	FullName          string    `json:"full_name" gorm:"column:full_name"`
	CreatedBy         int       `json:"created_by" gorm:"column:created_by"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedBy         int       `json:"updated_by" gorm:"column:updated_by"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at"`
	Active            bool      `json:"active" gorm:"column:active"`
}

func (DHLResPartnerIndustry) TableName() string {
	return "dhl_res_partner_industry"
}
