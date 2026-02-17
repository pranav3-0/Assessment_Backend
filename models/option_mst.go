package models

type OptionMst struct {
	OptionID    int64 `gorm:"column:option_id;primaryKey;autoIncrement"`
	ContentID   int64 `gorm:"column:content_id"`
	IsAnswer    bool  `gorm:"column:is_answer"`
	QuestionID  int64 `gorm:"column:question_id"`
	SequenceID  int64 `gorm:"column:sequence_id"`
	AnswerScore int   `gorm:"column:answer_score"`
}

func (OptionMst) TableName() string {
	return "option_mst"
}

type OptionMain struct {
	OptionID    int64  `gorm:"column:option_id" json:"option_id"`
	QuestionID  int64  `gorm:"column:question_id" json:"question_id"`
	OptionLabel string `gorm:"column:value" json:"option_label"`
	SequenceID  int64  `gorm:"column:sequence_id" json:"sequence_id"`
	IsAnswer    bool   `gorm:"column:is_answer" json:"is_answer"`
	AnswerScore int    `gorm:"column:answer_score" json:"answer_score"`
}
