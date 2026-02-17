package models

import "time"

type QuestionMst struct {
	QuestionID     int64     `gorm:"column:question_id;primaryKey;autoIncrement"`
	CreatedOn      time.Time `gorm:"column:created_on"`
	CreatedBy      string    `gorm:"column:created_by"`
	IsActive       bool      `gorm:"column:is_active"`
	IsDeleted      bool      `gorm:"column:is_deleted"`
	ModifiedOn     time.Time `gorm:"column:modified_on"`
	ModifiedBy     string    `gorm:"column:modified_by"`
	ContentID      int64     `gorm:"column:content_id"`
	QuestionTypeID int64     `gorm:"column:question_type_id"`
}

func (QuestionMst) TableName() string {
	return "question_mst"
}

type QuestionMain struct {
	QuestionID    int64  `gorm:"column:question_id" json:"question_id"`
	QuestionType  string `gorm:"column:type" json:"question_type"`
	QuestionTitle string `gorm:"column:value" json:"title"`
}

type QuestionContentWithType struct {
	ContentID      int64  `gorm:"column:content_id"`
	ContentType    string `gorm:"column:content_type"`
	Font           string `gorm:"column:font"`
	Value          string `gorm:"column:value"`
	QuestionTypeId uint64 `gorm:"column:question_type_id"`
	QuestionType   string `gorm:"column:question_type"`
}
