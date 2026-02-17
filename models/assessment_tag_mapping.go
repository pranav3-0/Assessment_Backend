package models

import "time"

type AssessmentTagMapping struct {
	AssessmentTagID    int64     `gorm:"column:assessment_tag_id;primaryKey;autoIncrement"`
	AssessmentSequence string    `gorm:"column:assessment_sequence"`
	TagID              int64     `gorm:"column:tag_id"`
	CreatedOn          time.Time `gorm:"column:created_on"`
	CreatedBy          string    `gorm:"column:created_by"`
	IsActive           bool      `gorm:"column:is_active"`
	IsDeleted          bool      `gorm:"column:is_deleted"`
	ModifiedOn         time.Time `gorm:"column:modified_on"`
	ModifiedBy         string    `gorm:"column:modified_by"`
}

func (AssessmentTagMapping) TableName() string {
	return "assessment_tag_mapping"
}
