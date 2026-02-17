package models

import (
	"time"

	"github.com/google/uuid"
)

type AssessmentUserSession struct {
	SessionID      uuid.UUID `gorm:"column:session_id;type:uuid;primaryKey" json:"assessment_session_id"`
	CreatedOn      time.Time `gorm:"column:created_on" json:"created_on"`
	CreatedBy      string    `gorm:"column:created_by" json:"created_by"`
	IsActive       bool      `gorm:"column:is_active" json:"is_active"`
	IsDeleted      bool      `gorm:"column:is_deleted" json:"is_deleted"`
	ModifiedOn     time.Time `gorm:"column:modified_on" json:"modified_on"`
	ModifiedBy     string    `gorm:"column:modified_by" json:"modified_by"`
	AccessTime     time.Time `gorm:"column:access_time" json:"access_time"`
	PartnerID      int64     `gorm:"column:partner_id" json:"partner_id"`
	UserID         string    `gorm:"column:user_id" json:"user_id"`
	AssessmentID   string    `gorm:"column:assessment_id" json:"assessment_sequence"`
	AssessmentType string    `gorm:"column:assessment_type" json:"assessment_type"`
	OdooToken      string    `gorm:"column:odoo_token" json:"-"`
}

func (AssessmentUserSession) TableName() string {
	return "assessment_user_session"
}
