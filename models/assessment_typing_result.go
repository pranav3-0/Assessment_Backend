package models

import "time"

type AssessmentTypingResult struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	AssessmentSeq   string    `gorm:"column:assessment_sequence;not null" json:"assessment_sequence" binding:"required"`
	WPM             int64     `gorm:"column:wpm" json:"wpm"`
	Accuracy        float64   `gorm:"column:accuracy" json:"accuracy"`
	TotalWords      int64     `gorm:"column:total_words" json:"total_words"`
	CorrectWords    int64     `gorm:"column:correct_words" json:"correct_words"`
	IncorrectWords  int64     `gorm:"column:incorrect_words" json:"incorrect_words"`
	CharactersTyped int64     `gorm:"column:characters_typed" json:"characters_typed"`
	UserID          string    `gorm:"column:user_id;not null" json:"user_id"`
	CreatedOn       time.Time `gorm:"column:created_on;autoCreateTime" json:"created_on"`
}

func (AssessmentTypingResult) TableName() string {
	return "assessment_typing_result"
}
