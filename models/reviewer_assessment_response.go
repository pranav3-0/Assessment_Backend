package models

type ReviewerAssessmentResponse struct {
	AssessmentID       int      `json:"assessment_id"`
	AssessmentSequence string   `json:"assessment_sequence"`
	AssessmentTitle    string   `json:"assessment_title"`

	AuthorName         string   `json:"author_name"`
	JobTitle           *string  `json:"job_title"`

	Subject            *string  `json:"subject"`
	Topics             []string `json:"topics"`

	Marks              int      `json:"marks"`
	TimeLimit          float64  `json:"time_limit"`
	Deadline           *string  `json:"deadline"`

	ReviewStatus       string   `json:"review_status"`
}