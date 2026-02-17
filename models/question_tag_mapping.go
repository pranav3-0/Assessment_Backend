package models

import "time"

type QuestionTagMapping struct {
	QuestionTagID int64     `gorm:"column:question_tag_id;primaryKey;autoIncrement"`
	QuestionID    int64     `gorm:"column:question_id"`
	TagID         int64     `gorm:"column:tag_id"`
	CreatedOn     time.Time `gorm:"column:created_on"`
	CreatedBy     string    `gorm:"column:created_by"`
	IsActive      bool      `gorm:"column:is_active"`
	IsDeleted     bool      `gorm:"column:is_deleted"`
	ModifiedOn    time.Time `gorm:"column:modified_on"`
	ModifiedBy    string    `gorm:"column:modified_by"`
}


func (QuestionTagMapping) TableName() string {
	return "tag_question_mapping"
}
