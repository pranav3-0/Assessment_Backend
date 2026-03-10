package models


type ReviewAssessmentRequest struct {
	AssessmentSequence string  `json:"assessment_sequence" binding:"required"`
	Status             string  `json:"status" binding:"required"` // approved / rejected
	Note               *string `json:"note"`
}