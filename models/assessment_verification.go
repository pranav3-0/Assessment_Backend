package models

import "time"

type AssessmentVerification struct {
    VerificationID     int       `json:"verification_id" gorm:"column:verification_id"`
    AssessmentSequence string    `json:"assessment_sequence" gorm:"column:assessment_sequence"`
    UserID             string    `json:"user_id" gorm:"column:user_id"`
    SessionID          *string   `json:"session_id" gorm:"column:session_id"`
    Photo              []byte    `json:"-" gorm:"column:photo"`
    Voice              []byte    `json:"-" gorm:"column:voice"`
    CreatedOn          time.Time `json:"created_on" gorm:"column:created_on"`
    UpdatedOn          time.Time `json:"updated_on" gorm:"column:updated_on"`
}
