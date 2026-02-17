package models

type ContentMst struct {
	ContentID     int64  `gorm:"column:content_id;primaryKey;autoIncrement"`
	ContentTypeID int64  `gorm:"column:content_type_id"`
	Font          string `gorm:"column:font"`
	Value         string `gorm:"column:value"`
}

func (ContentMst) TableName() string {
	return "content_mst"
}

type ContentWithType struct {
	ContentID   int64  `gorm:"column:content_id"`
	ContentType string `gorm:"column:content_type"`
	Font        string `gorm:"column:font"`
	Value       string `gorm:"column:value"`
}
