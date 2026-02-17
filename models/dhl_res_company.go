package models

import "time"

type DHLResCompany struct {
	CompanyID              int64     `json:"company_id" gorm:"primaryKey;column:company_id"`
	Name                   string    `json:"name" gorm:"column:name"`
	PartnerID              int       `json:"partner_id" gorm:"column:partner_id"`
	CurrencyID             int       `json:"currency_id" gorm:"column:currency_id"`
	Sequence               int       `json:"sequence" gorm:"column:sequence"`
	CreatedAt              time.Time `json:"created_at" gorm:"column:created_at"`
	ParentID               int       `json:"parent_id" gorm:"column:parent_id"`
	ReportHeader           string    `json:"report_header" gorm:"column:report_header"`
	ReportFooter           string    `json:"report_footer" gorm:"column:report_footer"`
	LogoWeb                []byte    `json:"logo_web" gorm:"column:logo_web"` // bytea
	AccountNo              string    `json:"account_no" gorm:"column:account_no"`
	Email                  string    `json:"email" gorm:"column:email"`
	Phone                  string    `json:"phone" gorm:"column:phone"`
	CompanyRegistry        string    `json:"company_registry" gorm:"column:company_registry"`
	PaperFormatID          int       `json:"paperformat_id" gorm:"column:paperformat_id"`
	ExternalReportLayoutID int       `json:"external_report_layout_id" gorm:"column:external_report_layout_id"`
	BaseOnboardingState    string    `json:"base_onboarding_company_state" gorm:"column:base_onboarding_company_state"`
	Font                   string    `json:"font" gorm:"column:font"`
	PrimaryColor           string    `json:"primary_color" gorm:"column:primary_color"`
	SecondaryColor         string    `json:"secondary_color" gorm:"column:secondary_color"`
	CreatedBy              int       `json:"created_by" gorm:"column:created_by"`
	UpdatedBy              int       `json:"updated_by" gorm:"column:updated_by"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"column:updated_at"`
	SocialTwitter          string    `json:"social_twitter" gorm:"column:social_twitter"`
	SocialFacebook         string    `json:"social_facebook" gorm:"column:social_facebook"`
	SocialGithub           string    `json:"social_github" gorm:"column:social_github"`
	SocialLinkedIn         string    `json:"social_linkedin" gorm:"column:social_linkedin"`
	SocialYoutube          string    `json:"social_youtube" gorm:"column:social_youtube"`
	SocialInstagram        string    `json:"social_instagram" gorm:"column:social_instagram"`
	PartnerGID             int       `json:"partner_gid" gorm:"column:partner_gid"`
	SnailmailColor         bool      `json:"snailmail_color" gorm:"column:snailmail_color"`
	SnailmailCover         bool      `json:"snailmail_cover" gorm:"column:snailmail_cover"`
	SnailmailDuplex        bool      `json:"snailmail_duplex" gorm:"column:snailmail_duplex"`
	IsActive               bool      `json:"is_active" gorm:"column:is_active"`
}

func (DHLResCompany) TableName() string {
	return "dhl_res_company"
}
