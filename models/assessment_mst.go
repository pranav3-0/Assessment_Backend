package models

import (
	"time"

	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

type AssessmentMst struct {
	AssessmentID       int64      `gorm:"column:assessment_id;primaryKey;autoIncrement"`
	CreatedOn          time.Time  `gorm:"column:created_on"`
	CreatedBy          string     `gorm:"column:created_by"`
	IsActive           bool       `gorm:"column:is_active"`
	IsDeleted          bool       `gorm:"column:is_deleted"`
	ModifiedOn         time.Time  `gorm:"column:modified_on"`
	ModifiedBy         string     `gorm:"column:modified_by"`
	AssessmentDesc     string     `gorm:"column:assessment_desc"`
	AssessmentSequence string     `gorm:"column:assessment_sequence"`
	Duration           int64      `gorm:"column:duration"`
	Marks              int64      `gorm:"column:marks"`
	StartTime          *time.Time `gorm:"column:start_time"`
	PartnerID          int64      `gorm:"column:partner_id"`
	ValidFrom          *time.Time `gorm:"column:valid_from"`
	ValidTo            *time.Time `gorm:"column:valid_to"`
	NoFixedSchedule    bool       `gorm:"column:no_fixed_schedule"`
	Instructions       string     `gorm:"column:instructions"`
	JobID              *int64     `gorm:"column:job_id"`
	AssessmentType     string     `gorm:"column:assessment_type"`
}

func (AssessmentMst) TableName() string {
	return "assessment_mst"
}
