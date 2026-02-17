package models

import "time"

type ManagerAssessmentMapping struct {
	ID                 int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ManagerID          string    `gorm:"column:manager_id"`
	AssessmentSequence string    `gorm:"column:assessment_sequence"`
	IsActive           bool      `gorm:"column:is_active"`
	CreatedOn          time.Time `gorm:"column:created_on"`
}

func (ManagerAssessmentMapping) TableName() string {
	return "manager_assessment_mapping"
}
