package models

import "time"

type AssessmentReportFilter struct {
	FromDate           *time.Time `json:"from_date"`
	ToDate             *time.Time `json:"to_date"`
	AssessmentID       string     `json:"assessment_id"`
	CenterID           *int64     `json:"center_id"`
	CompanyID          *int64     `json:"company_id"`
	EmployeeName       string     `json:"employee_name"`
	Status             string     `json:"status"`
	SkillSet           string     `json:"skill_set"`
	QuestionnaireTitle string     `json:"questionnaire_title"`
}

type AssessmentReportRow struct {
	EmployeeName    string    `json:"employee_name"`
	TeamLead        string    `json:"team_lead"`
	TeamManager     string    `json:"team_manager"`
	SeniorManager   string    `json:"senior_manager"`
	SDL             string    `json:"sdl"`
	SLL             string    `json:"sll"`
	SkillSet        string    `json:"skill_set"`
	AssessmentDate  time.Time `json:"assessment_date"`
	AssessmentTitle string    `json:"assessment_title"`
	Status          string    `json:"status"`
	Attempts        int       `json:"attempts"`

	TotalAssigned   int `json:"total_assigned"`
	TotalPassed     int `json:"total_passed"`
	TotalFailed     int `json:"total_failed"`
	TotalNotStarted int `json:"total_not_started"`

	AssessmentSequence string `json:"assessment_sequence"`
}

type AssessmentReportResponse struct {
	Filters AssessmentReportFilter `json:"filters"`
	Rows    []AssessmentReportRow  `json:"rows"`
	Total   int                    `json:"total"`
}
