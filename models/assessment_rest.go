package models

import (
	"time"
	"github.com/lib/pq"
)

type GetAssessmentRequest struct {
    AssessmentId string `json:"assessment_sequence"`
    SessionID    string `json:"session_id"`
    UserId       string `json:"-"`
}

type DuplicateAssessmentRequest struct {
	AssessmentSequence string `json:"assessment_sequence"`
}

type AssessmentResponse struct {
	AssessmentID          int64                `json:"assessment_id"`
	AssessmentSequence    string               `json:"assessment_sequence,omitempty"`
	AssessmentName        string               `json:"assessment_name"`
	AssessmentUsersStatus string               `json:"assessment_user_status"`
	AssessmentStatus      string               `json:"assessment_status"`
	QuestionsCount        int                  `json:"questions_count"`
	AssessmentType        string               `json:"assessment_type"`
	AssessmentDuration    *int64               `json:"assessment_duration"`
	NoOfAttempts          int                  `json:"no_of_attempts"`
	TimeLimit             float64              `json:"time_limit"`
	Marks                 int64                `json:"marks"`
	Deadline              *time.Time           `json:"deadline"`
	CenterName            *string              `json:"center_name"`
	ServiceLineName       *string              `json:"service_line_name"`
	BusinessPartnerName   *string              `json:"business_partner_name"`
	SubBusinessPartner    *string              `json:"sub_business_partner_name"`
	ServiceGroupName      *string              `json:"service_group_name"`
	ServiceName           *string              `json:"service_name"`
	SessionID             string               `json:"sessionId,omitempty"`
	AssessmentReport      bool                 `json:"assessment_report,omitempty"`
	Certificate           bool                 `json:"certificate"`
	AttemptedTypingTest   bool                 `json:"attempted_typing_test"`
	IsTypingTest          bool                 `json:"is_typing_test"`
	Instruction           string               `json:"instruction"`
	Questions             []AssessmentQuestion `json:"questions"`
	SessionImage          []byte               `json:"session_image,omitempty"`
	Tags                  []TagRequest         `json:"tags,omitempty"`
}

type AssessmentQuestion struct {
	QuestionID      int          `json:"question_id"`
	Sequence        int          `json:"sequence"`
	Title           string       `json:"title"`
	Answers         []Answer     `json:"answers"`
	AttemptedAnswer *int         `json:"attempted_answer"`
	QuestionTime    int          `json:"question_time"`
	QuestionTypeId  uint64       `json:"question_type_id"`
	QuestionType    string       `json:"question_type"`
	SkippingAllowed bool         `json:"skipping_allowed"`
	Tags            []TagRequest `json:"tags,omitempty"`
}

type Answer struct {
	AnswerID      int    `json:"option_id"`
	Sequence      *int64 `json:"sequence"`
	OptionLabel   string `json:"option_label"`
	CorrectAnswer bool   `json:"correctAnswer,omitempty"`
}

type AssessmentFilter struct {
	AssessmentSequence  *string `json:"assessment_sequence"`
	AssessmentSessionId *string `json:"assessment_session"`
	AssessmentID        *string `json:"assessment_id"`
}

type AssessmentListResponse struct {
	AssessmentID        int                     `json:"assessment_id" gorm:"column:assessment_id"`
	AssessmentSequence  string                  `json:"assessment_sequence" gorm:"column:assessment_sequence"`
	AssessmentTitle     string                  `json:"assessment_title" gorm:"column:assessment_title"`
	TimeLimit           *int                    `json:"time_limit" gorm:"column:time_limit"`
	AssessmentType      *string                 `json:"assessment_type"`
	Marks               *float64                `json:"marks" gorm:"column:marks"`
	PassingScore        *float64                `json:"passing_score,omitempty" gorm:"column:passing_score"`
	MarksObtained       *float64                `json:"marks_obtained,omitempty" gorm:"column:marks_obtained"`
	UserScore           *float64                `json:"user_score,omitempty" gorm:"column:user_score"`
	Deadline            *time.Time              `json:"deadline" gorm:"column:deadline"`
	State               *string                 `json:"state" gorm:"column:state"`
	CenterName          *string                 `json:"center_name" gorm:"column:center_name"`
	ServiceLineName     *string                 `json:"service_line_name" gorm:"column:service_line_name"`
	BusinessPartnerName *string                 `json:"business_partner_name" gorm:"column:business_partner_name"`
	SubBusinessPartner  *string                 `json:"sub_business_partner_name" gorm:"column:sub_business_partner_name"`
	ServiceGroupName    *string                 `json:"service_group_name" gorm:"column:service_group_name"`
	ServiceName         *string                 `json:"service_name" gorm:"column:service_name"`
	ShowResult          bool                    `json:"show_result" gorm:"column:show_result"`
	ResultStatus        string                  `json:"result_status,omitempty" gorm:"column:result_status"`
	Certificate         bool                    `json:"certificate" gorm:"column:certificate"`
	AttemptedTypingTest bool                    `json:"attempted_typing_test" gorm:"column:attempted_typing_test"`
	IsTypingTest        bool                    `json:"is_typing_test" gorm:"column:is_typing_test"`
	TypingResult        *AssessmentTypingResult `json:"typing_result" gorm:"-"`
	SessionImage        []byte                  `json:"session_image,omitempty" gorm:"-"`
	JobTitle            *string                 `json:"job_title"`
	Tags pq.StringArray `json:"tags" gorm:"column:tags;type:text[]"`
}

type AssessmentAttendeesInfo struct {
	AssessmentID     string  `json:"assessmentId" gorm:"column:assessment_id"`
	SessionID        string  `json:"sessionId" gorm:"column:session_id"`
	UserID           string  `json:"userId" gorm:"column:user_id"`
	AssessmentStatus string  `json:"assessmentStatus" gorm:"column:assessment_status"`
	FirstName        string  `json:"firstName" gorm:"column:first_name"`
	LastName         string  `json:"lastName" gorm:"column:last_name"`
	Email            string  `json:"email" gorm:"column:email"`
	CenterName       *string `json:"centerName" gorm:"column:center_name"`
	TeamLead         *string `json:"teamLead" gorm:"column:team_lead"`
	Manager          *string `json:"manager" gorm:"column:manager"`
	SeniorManager    *string `json:"seniorManager" gorm:"column:senior_manager"`
	SDL              *string `json:"sdl" gorm:"column:sdl"`
	SLL              *string `json:"sll" gorm:"column:sll"`
}

type AssessmentUserResponse struct {
	QuestionID         int64    `json:"question_id" gorm:"column:question_id"`
	QuestionText       string   `json:"question_text" gorm:"column:question_text"`
	SelectedOptionID   *int64   `json:"selected_option_id" gorm:"column:selected_option_id"`
	SelectedOptionText *string  `json:"selected_option_text" gorm:"column:selected_option_text"`
	PointAssigned      *float64 `json:"point_assigned" gorm:"column:point_assigned"`
	AnswerStatus       string   `json:"answer_status" gorm:"column:answer_status"`
	Answers            []Answer `json:"options" gorm:"-"`
	// CorrectOptionTexts pq.StringArray `json:"correct_option_texts" gorm:"column:correct_option_texts"`
	// CorrectOptionIDs   pq.Int64Array  `json:"correct_option_ids" gorm:"column:correct_option_ids"`
}

// Upload Assessment
type SheetOption struct {
	Label     string `json:"label"`
	IsCorrect bool   `json:"is_correct"`
	Score     int    `json:"score"`
}

// TagRequest represents a tag with optional parent and children
type TagRequest struct {
	ParentTag string   `json:"parent_tag,omitempty"`
	ChildTags []string `json:"child_tags,omitempty"`
}

type SheetQuestion struct {
	Title             string        `json:"title"`
	MandatoryToAnswer bool          `json:"mandatory_to_answer"`
	QuestionType      string        `json:"question_type"`
	Options           []SheetOption `json:"options"`
	Tags              []TagRequest  `json:"tags,omitempty"`
}

type SheetAssessment struct {
	AssessmentName     string          `json:"assessment_name"`
	AssessmentSequence string          `json:"assessment_sequence"`
	Questions          []SheetQuestion `json:"questions"`
	Tags               []TagRequest    `json:"tags,omitempty"`
}

// Update Assessment Data Models
type UpdateAssessmentRequest struct {
	AssessmentSequence string           `json:"assessment_sequence"`
	AssessmentDetails  AssessmentUpdate `json:"assessment_details"`
	Questions          []QuestionDTO    `json:"questions"`
	DeletedQuestionIDs []int64          `json:"deleted_question_ids"`
	Tags               []TagRequest     `json:"tags,omitempty"` 
}

type AssessmentUpdate struct {
	AssessmentDesc       *string    `json:"assessment_desc"`
	AssessmentType       *string    `json:"assessment_type"`
	Duration             *int64     `json:"duration"`
	Marks                *int64     `json:"marks"`
	State                *string    `json:"state"`
	TimeLimit            *int       `json:"time_limit"`
	CenterId             *int       `json:"center_id"`
	ServiceLineID        *int       `json:"service_line_id"`
	BusinessPartnerID    *int       `json:"business_partner_id"`
	SubBusinessPartnerID *int       `json:"sub_business_partner_id"`
	ServiceGroupID       *int       `json:"service_group_id"`
	ServiceID            *int       `json:"service_id"`
	Deadline             *time.Time `json:"deadline"`
	Instruction          *string    `json:"instruction"`
	AllowShowResult      *bool      `json:"allow_show_result"`
	AllowViewCertificate *bool      `json:"allow_view_certificat"`
	JobID                *int64     `json:"job_id"`
}

type QuestionDTO struct {
	QuestionID   int64        `json:"question_id,omitempty"`
	Title        string       `json:"title"`
	QuestionType string       `json:"question_type"`
	Options      []OptionDTO  `json:"options"`
	Tags         []TagRequest `json:"tags,omitempty"`
}

type OptionDTO struct {
	OptionID  int64  `json:"option_id,omitempty"`
	Label     string `json:"label"`
	IsCorrect bool   `json:"is_correct"`
	Score     int    `json:"score"`
}

// Update Assessment Status
type UpdateAssessmentStatusRequest struct {
	AssessmentSequence string `json:"assessment_sequence"  binding:"required"`
	AssessmentStatus   string `json:"assessment_status"  binding:"required"`
}

type SubmitUserAssessmentRequest struct {
	AssessmentSequence string         `json:"assessment_sequence" binding:"required"`
	SessionID          string         `json:"session_id" binding:"required"`
	Response           []UserResponse `json:"response"`
}

type UserResponse struct {
	QuestionID       int64   `json:"question_id"`
	SelectedOptionID []int64 `json:"selected_option_id"`
}

type DistributeAssessmentRequest struct {
	AssessmentSequence string   `json:"assessment_sequence" binding:"required"`
	UserIds            []string `json:"user_ids" binding:"required"`
}

type MapUsersToManagerRequest struct {
	ManagerID string   `json:"manager_id" binding:"required"`
	UserIDs   []string `json:"user_ids" binding:"required"`
}

type CertificateDetails struct {
	UserID             string
	UserName           string
	AssessmentSequence string
	AssessmentTitle    string
	TotalMarks         int
	MarksObtained      int
	PassingScore       float64
	SessionID          string
	CompletedAt        time.Time
}

type ManualAssessmentRequest struct {
	AssessmentName       string       `json:"assessment_name"`
	Duration             int64        `json:"duration"`
	Marks                int64        `json:"marks"`
	StartTime            *time.Time   `json:"start_time"`
	ValidFrom            *time.Time   `json:"valid_from"`
	ValidTo              *time.Time   `json:"valid_to"`
	Instructions         string       `json:"instructions"`
	AssessmentType       string       `json:"assessment_type"`
	Certificate          bool         `json:"certificate"`
	CenterID             int          `json:"center_id"`
	ServiceLineID        int          `json:"service_line_id"`
	BusinessPartnerID    int          `json:"business_partner_id"`
	SubBusinessPartnerID int          `json:"sub_business_partner_id"`
	ServiceGroupID       int          `json:"service_group_id"`
	ServiceID            int          `json:"service_id"`
	Questions            []int64      `json:"questions"`
	JobID                *int64       `json:"job_id"`
	Tags                 []TagRequest `json:"tags,omitempty"`
}

type GenerateAssessmentRequest struct {
	Topic             string       `json:"topic" binding:"required"`
	NumberOfQuestions int          `json:"number_of_questions" binding:"required,min=1,max=50"`
	DifficultyLevel   int          `json:"difficulty_level" binding:"required,min=1,max=5"`
	QuestionType      string       `json:"question_type"`
	Tags              []TagRequest `json:"tags,omitempty"`
	AssessmentName    string       `json:"assessment_name"`
}

type GenerateAssessmentResponse struct {
	Assessment SheetAssessment `json:"assessment"`
}

type SaveGeneratedAssessmentRequest struct {
	Assessment SheetAssessment `json:"assessment" binding:"required"`
}

type StartAssessmentRequest struct {
	AssessmentSequence string `json:"assessment_sequence" binding:"required"`
}

type StartAssessmentResponse struct {
	SessionID          string `json:"session_id"`
	AssessmentSequence string `json:"assessment_sequence"`
}
