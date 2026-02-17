package models

import "time"

type ContactUsResponse struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Company     string    `json:"company"`
	Subject     string    `json:"subject"`
	Question    string    `json:"question"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ContactUsResponse) TableName() string {
	return "dhl_contact_us_responses"
}
