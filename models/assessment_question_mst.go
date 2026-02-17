package models

import "time"

type AssessmentQuestionMst struct {
	AssessmentQuestionID int64     `gorm:"column:assessment_question_id;primaryKey;autoIncrement"`
	CreatedOn            time.Time `gorm:"column:created_on"`
	CreatedBy            string    `gorm:"column:created_by"`
	IsActive             bool      `gorm:"column:is_active"`
	IsDeleted            bool      `gorm:"column:is_deleted"`
	ModifiedOn           time.Time `gorm:"column:modified_on"`
	ModifiedBy           string    `gorm:"column:modified_by"`
	AssessmentSequence   string    `gorm:"column:assessment_sequence"`
	CorrectPoints        int64     `gorm:"column:correct_points"`
	DurationInSeconds    int64     `gorm:"column:duration_in_seconds"`
	NegativePoints       int64     `gorm:"column:negative_points"`
	QuestionID           int64     `gorm:"column:question_id"`
	SequenceID           int64     `gorm:"column:sequence_id"`
	SkippingAllowed      bool      `gorm:"column:skipping_allowed"`
	DifficultyLevel      string    `gorm:"column:difficulty_level"`
	AssessmentID         int       `gorm:"column:assessment_id"`
}

func (AssessmentQuestionMst) TableName() string {
	return "assessment_question_mst"
}
