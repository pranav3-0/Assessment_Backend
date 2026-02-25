package models

type AdminAssessmentUserResultResponse struct {
	UserID             string                        `json:"user_id"`
	AssessmentSequence string                        `json:"assessment_sequence"`
	Attempts           []AssessmentAttemptResult     `json:"attempts"`
}

type AssessmentAttemptResult struct {
	SessionID        string                    `json:"session_id"`
	AttemptTime      interface{}               `json:"attempt_time"`
	TotalMarks       int64                     `json:"total_marks"`
	ObtainedMarks    int64                     `json:"obtained_marks"`
	CompletionStatus string                    `json:"completion_status"`
	Questions        []AdminQuestionResult     `json:"questions"`
}

type AdminQuestionResult struct {
	QuestionID        int64  `json:"question_id"`
	QuestionText      string `json:"question_text"`
	SelectedOption    string `json:"selected_option"`
	CorrectOption     string `json:"correct_option"`
	IsCorrect         bool   `json:"is_correct"`
	PointsAssigned    int64  `json:"points_assigned"`
	CorrectPoints     int64  `json:"correct_points"`
	NegativePoints    int64  `json:"negative_points"`
	AttemptCount      int64  `json:"attempt_count"`
	AttemptTime       int64  `json:"attempt_time"`
}