package models

import "time"

type JobDescription struct {
	JobID          int64     `gorm:"column:job_id;primaryKey;autoIncrement" json:"job_id"`
	CreatedOn      time.Time `gorm:"column:created_on" json:"created_on"`
	CreatedBy      string    `gorm:"column:created_by" json:"created_by"`
	IsActive       bool      `gorm:"column:is_active" json:"is_active"`
	IsDeleted      bool      `gorm:"column:is_deleted" json:"is_deleted"`
	ModifiedOn     time.Time `gorm:"column:modified_on" json:"modified_on"`
	ModifiedBy     string    `gorm:"column:modified_by" json:"modified_by"`
	Title          string    `gorm:"column:title" json:"title"`
	Description    string    `gorm:"column:description" json:"description"`
	RequiredSkills string    `gorm:"column:required_skills" json:"required_skills"`
	Level          string    `gorm:"column:level" json:"level"`
}


func (JobDescription) TableName() string {
	return "job_descriptions"
}
