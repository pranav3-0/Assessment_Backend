package models

import "time"

type AssessmentResult struct {
	AssessmentResultID  int64      `gorm:"column:assessment_result_id;primaryKey;autoIncrement"`
	CreatedOn           time.Time  `gorm:"column:created_on"`
	CreatedBy           string     `gorm:"column:created_by"`
	IsActive            bool       `gorm:"column:is_active"`
	IsDeleted           bool       `gorm:"column:is_deleted"`
	ModifiedOn          time.Time  `gorm:"column:modified_on"`
	ModifiedBy          string     `gorm:"column:modified_by"`
	ActivityID          int64      `gorm:"column:activity_id"`
	AssessmentSequence  string     `gorm:"column:assessment_sequence"`
	AssessmentSessionID string     `gorm:"column:assessment_session_id"`
	AttemptEndTime      *time.Time `gorm:"column:attempt_end_time"`
	AttemptID           int64      `gorm:"column:attempt_id"`
	AttemptOptionID     int64      `gorm:"column:attempt_option_id"`
	AttemptStartTime    *time.Time `gorm:"column:attempt_start_time"`
	Description         string     `gorm:"column:description"`
	PointAssigned       int64      `gorm:"column:point_assigned"`
	QuestionID          int64      `gorm:"column:question_id"`
	AttemptCount        int64      `gorm:"column:attempt_count"`
	AttemptTime         int64      `gorm:"column:attempt_time"`
	BonusPointAssigned  int64      `gorm:"column:bonus_point_assigned"`
	SelectedOptionIDs   string     `gorm:"column:selected_option_ids"`
}

func (AssessmentResult) TableName() string {
	return "assessment_result"
}
