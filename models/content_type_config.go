package models

type ContentTypeConfig struct {
	ContentTypeID int64  `gorm:"column:content_type_id;primaryKey;autoIncrement"`
	ContentType   string `gorm:"column:content_type"`
}

func (ContentTypeConfig) TableName() string {
	return "content_type_config"
}
