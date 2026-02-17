package models

import "time"

type TagMaster struct {
	TagID       int64     `gorm:"column:tag_id;primaryKey;autoIncrement"`
	TagName     string    `gorm:"column:tag"`
	ParentTagID *int64    `gorm:"column:parent_tag_id"`
	CreatedOn   time.Time `gorm:"column:created_on"`
	CreatedBy   string    `gorm:"column:created_by"`
	IsActive    bool      `gorm:"column:is_active"`
	IsDeleted   bool      `gorm:"column:is_deleted"`
	ModifiedOn  time.Time `gorm:"column:modified_on"`
	ModifiedBy  string    `gorm:"column:modified_by"`
}

func (TagMaster) TableName() string {
	return "tag_mst"
}
