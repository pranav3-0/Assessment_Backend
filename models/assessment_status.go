package models

import "time"

type AssessmentStatus struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedOn        time.Time `gorm:"column:created_on"`
	CreatedBy        string    `gorm:"column:created_by"`
	IsActive         bool      `gorm:"column:is_active"`
	IsDeleted        bool      `gorm:"column:is_deleted"`
	ModifiedOn       time.Time `gorm:"column:modified_on"`
	ModifiedBy       string    `gorm:"column:modified_by"`
	AssessmentID     string    `gorm:"column:assessment_id"`
	AssessmentStatus string    `gorm:"column:assessment_status"`
	UserID           string    `gorm:"column:user_id"`
}

func (AssessmentStatus) TableName() string {
	return "assessment_status"
}
