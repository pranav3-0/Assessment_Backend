package models

import (
	"time"

	"github.com/google/uuid"
)

type AssessmentUserSessionImage struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SessionID     uuid.UUID `gorm:"column:session_id;type:uuid;not null" json:"session_id"`
	Image         []byte    `gorm:"column:image;type:bytea;not null" json:"-"`
	CreatedOn     time.Time `gorm:"column:created_on;not null" json:"created_on"`
	CreatedBy     string    `gorm:"column:created_by" json:"created_by"`
	IsActive      bool      `gorm:"column:is_active;default:true" json:"is_active"`
	IsDeleted     bool      `gorm:"column:is_deleted;default:false" json:"is_deleted"`
	ModifiedOn    time.Time `gorm:"column:modified_on" json:"modified_on"`
	ModifiedBy    string    `gorm:"column:modified_by" json:"modified_by"`
}

func (AssessmentUserSessionImage) TableName() string {
	return "assessment_user_session_image"
}

// Request DTO for creating session image
type CreateSessionImageRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

// Response DTO for session image
type SessionImageResponse struct {
	ID        int64     `json:"id"`
	SessionID uuid.UUID `json:"session_id"`
	CreatedOn time.Time `json:"created_on"`
}
