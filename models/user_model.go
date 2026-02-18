package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AssessmentUser struct {
	UserID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:user_id" json:"user_id"`
	FirstName  string    `gorm:"type:varchar(255);column:first_name" json:"first_name"`
	LastName   string    `gorm:"type:varchar(255);column:last_name" json:"last_name"`
	Email      string    `gorm:"type:varchar(255);column:email" json:"email"`
	Phone      string    `gorm:"type:varchar(255);column:phone" json:"phone"`
	Username   string    `gorm:"type:varchar(255);not null;column:username" json:"username"`
	AuthUserID string    `gorm:"type:varchar(255);column:auth_user_id" json:"auth_user_id"`
	Password   string    `gorm:"type:varchar(255);column:password" json:"-"`
	IsActive   bool      `gorm:"default:true;column:is_active" json:"is_active"`
	NotifyId   string    `gorm:"column:notify_id;type:varchar" json:"notify_id"`

	UserType string `gorm:"type:varchar(20);column:user_type" json:"user_type"` // THIS

	CreatedAt time.Time  `gorm:"type:timestamptz;default:now();column:created_at" json:"created_at"`
	CreatedBy *uuid.UUID `gorm:"type:uuid;column:created_by" json:"created_by"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;default:now();column:updated_at" json:"updated_at"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid;column:updated_by" json:"updated_by"`
}

func (AssessmentUser) TableName() string {
	return "assessment_user_mst"
}

type DhlAssessmentUserMstExt struct {
	UserID           string    `gorm:"column:user_id"`
	CompanyID        *int      `gorm:"column:company_id"`
	Signature        *string   `gorm:"column:signature"`
	Share            *bool     `gorm:"column:share"`
	NotificationType *string   `gorm:"column:notification_type"`
	Karma            *int      `gorm:"column:karma"`
	RankID           *int      `gorm:"column:rank_id"`
	NextRankID       *int      `gorm:"column:next_rank_id"`
	TeamLead         *string   `gorm:"column:team_lead"`
	Manager          *string   `gorm:"column:manager"`
	SeniorManager    *string   `gorm:"column:senior_manager"`
	SDL              *string   `gorm:"column:sdl"`
	SLL              *string   `gorm:"column:sll"`
	Remark           *string   `gorm:"column:remark"`
	UserMap          *string   `gorm:"column:user_map"`
	EmpCode          *string   `gorm:"column:emp_code"`
	Center           *int      `gorm:"column:center"`
	IsUserCreated    *bool     `gorm:"column:is_user_created"`
	SelectionType    *string   `gorm:"column:selection_type"`
	ActionID         *int      `gorm:"column:action_id"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy        *string   `gorm:"column:updated_by"`
	CreatedBy        *string   `gorm:"column:created_by"`
}

func (DhlAssessmentUserMstExt) TableName() string {
	return "dhl_assessment_user_mst_ext"
}

type RegisterRequest struct {
	ClientName string   `json:"client_name" binding:"required"`
	FirstName  string   `json:"first_name" binding:"required"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email" binding:"required,email"`
	Phone      string   `json:"phone" binding:"required"`
	Roles      []string `json:"roles" binding:"required"`
	UserType   string   `json:"user_type"`
	Password   string   `json:"password"  binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"client_id"`
}

type LoginResponse struct {
	Token        string         `json:"token,omitempty"`
	RefreshToken string         `json:"refresh_token,omitempty"`
	User         *UserWithRoles `json:"user,omitempty"`
}

type UserWithRoles struct {
	UserID     uuid.UUID      `json:"user_id"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone"`
	Username   string         `json:"username"`
	AuthUserID string         `json:"auth_user_id"`
	Roles      pq.StringArray `json:"roles" gorm:"type:text[]"`
	UserType   string         `json:"user_type"` //  THIS
}

type UserProfileUpdate struct {
	FirstName     *string  `json:"first_name,omitempty"`
	LastName      *string  `json:"last_name,omitempty"`
	Email         *string  `json:"email,omitempty"`
	Phone         *string  `json:"phone,omitempty"`
	AuthUserID    *string  `json:"auth_user_id,omitempty"`
	NotifyId      *string  `json:"notify_id,omitempty"`
	Password      *string  `json:"password,omitempty"`
	CompanyID     *int     `json:"company_id,omitempty"`
	Karma         *int     `json:"karma,omitempty"`
	RankID        *int     `json:"rank_id,omitempty"`
	TeamLead      *string  `json:"team_lead,omitempty"`
	Manager       *string  `json:"manager,omitempty"`
	SeniorManager *string  `json:"senior_manager,omitempty"`
	SDL           *string  `json:"sdl,omitempty"`
	SLL           *string  `json:"sll,omitempty"`
	UserMap       *string  `json:"user_map,omitempty"`
	EmpCode       *string  `json:"emp_code,omitempty"`
	Center        *int     `json:"center,omitempty"`
	SelectionType *string  `json:"selection_type,omitempty"`
	Roles         []string `json:"roles,omitempty"`
}

type UserFullData struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Username  string    `json:"username"`

	CompanyID     *int     `json:"company_id"`
	Karma         *int     `json:"karma"`
	RankID        *int     `json:"rank_id"`
	TeamLead      *string  `json:"team_lead"`
	Manager       *string  `json:"manager"`
	SeniorManager *string  `json:"senior_manager"`
	SDL           *string  `json:"sdl"`
	SLL           *string  `json:"sll"`
	UserMap       *string  `json:"user_map"`
	EmpCode       *string  `json:"emp_code"`
	Center        *int     `json:"center"`
	SelectionType *string  `json:"selection_type"`
	Roles         []string `json:"roles"`
	UserType      string   `json:"user_type"`
}
