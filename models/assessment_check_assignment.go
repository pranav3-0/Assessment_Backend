package models

type CheckAssignmentRequest struct {
	AssessmentSequence string   `json:"assessment_sequence"`
	UserIDs            []string `json:"user_ids"`
}

type CheckAssignmentResponse struct {
	UserID            string `json:"user_id"`
	AssessmentStatus  string `json:"assessment_status"`
}