package models

import "time"

type DhlSurveySurveyExt struct {
	SurveyID                  int64     `gorm:"primaryKey;column:survey_id" json:"survey_id"`
	MessageMainAttachmentID   int       `gorm:"column:message_main_attachment_id" json:"message_main_attachment_id,omitempty"`
	Color                     int       `gorm:"column:color" json:"color,omitempty"`
	ThankYouMessage           string    `gorm:"column:thank_you_message" json:"thank_you_message,omitempty"`
	State                     string    `gorm:"column:state" json:"state,omitempty"`
	QuestionsLayout           string    `gorm:"column:questions_layout" json:"questions_layout,omitempty"`
	QuestionsSelection        string    `gorm:"column:questions_selection" json:"questions_selection,omitempty"`
	Category                  string    `gorm:"column:category" json:"category,omitempty"`
	AccessMode                string    `gorm:"column:access_mode" json:"access_mode,omitempty"`
	AccessToken               string    `gorm:"column:access_token" json:"access_token,omitempty"`
	UsersLoginRequired        bool      `gorm:"column:users_login_required" json:"users_login_required,omitempty"`
	UsersCanGoBack            bool      `gorm:"column:users_can_go_back" json:"users_can_go_back,omitempty"`
	ScoringType               string    `gorm:"column:scoring_type" json:"scoring_type,omitempty"`
	IsAttemptsLimited         bool      `gorm:"column:is_attempts_limited" json:"is_attempts_limited,omitempty"`
	AttemptsLimit             int       `gorm:"column:attempts_limit" json:"attempts_limit,omitempty"`
	IsTimeLimited             bool      `gorm:"column:is_time_limited" json:"is_time_limited,omitempty"`
	TimeLimit                 float64   `gorm:"column:time_limit" json:"time_limit,omitempty"`
	Certificate               bool      `gorm:"column:certificate" json:"certificate,omitempty"`
	CertificationMailTemplate int       `gorm:"column:certification_mail_template_id" json:"certification_mail_template_id,omitempty"`
	CertificationGiveBadge    bool      `gorm:"column:certification_give_badge" json:"certification_give_badge,omitempty"`
	CertificationBadgeID      int       `gorm:"column:certification_badge_id" json:"certification_badge_id,omitempty"`
	CenterID                  int       `gorm:"column:center_id" json:"center_id,omitempty"`
	ServiceLineID             int       `gorm:"column:service_line_id" json:"service_line_id,omitempty"`
	BusinessPartnerID         int       `gorm:"column:business_partner_id" json:"business_partner_id,omitempty"`
	SubBusinessPartnerID      int       `gorm:"column:sub_business_partner_id" json:"sub_business_partner_id,omitempty"`
	ServiceGroupID            int       `gorm:"column:service_group_id" json:"service_group_id,omitempty"`
	ServiceID                 int       `gorm:"column:service_id" json:"service_id,omitempty"`
	NoOfRandomQuestions       int       `gorm:"column:no_of_random_questions" json:"no_of_random_questions,omitempty"`
	ShowResult                bool      `gorm:"column:show_result" json:"show_result,omitempty"`
	IsTypingTest              bool      `gorm:"column:is_typing_test" json:"is_typing_test,omitempty"`
	TypingTestTime            float64   `gorm:"column:typing_test_time" json:"typing_test_time,omitempty"`
	SelectionType             string    `gorm:"column:selection_type" json:"selection_type,omitempty"`
	Deadline                  time.Time `gorm:"column:deadline" json:"deadline,omitempty"`
	AssessmentSequence        string    `gorm:"column:assessment_sequence" json:"assessment_sequence,omitempty"`
}

func (DhlSurveySurveyExt) TableName() string {
	return "dhl_survey_survey_ext"
}

type DhlSurveySurveyExtResponse struct {
	SurveyID               int64     `json:"survey_id"`
	State                  string    `json:"state"`
	TimeLimit              float64   `json:"time_limit"`
	Deadline               time.Time `json:"deadline"`
	AssessmentSequence     string    `json:"assessment_sequence"`
	AttemptsLimit          int       `json:"attempts_limit,omitempty"`
	CenterName             *string   `json:"center_name"`
	ServiceLineName        *string   `json:"service_line_name"`
	BusinessPartnerName    *string   `json:"business_partner_name"`
	SubBusinessPartnerName *string   `json:"sub_business_partner_name"`
	ServiceGroupName       *string   `json:"service_group_name"`
	ServiceName            *string   `json:"service_name"`
	Certificate            bool      `json:"certificate"`
	IsTypingTest           bool      `json:"is_typing_test"`
}
