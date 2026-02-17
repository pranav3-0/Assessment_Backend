package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

// SMTP Config
type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func SendEmail(cfg EmailConfig, to, subject, templateFile string, data interface{}) error {
	// Parse Template
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	// SMTP AUTH
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	// Compose Email Message
	msg := []byte(
		"From: " + cfg.From + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
			body.String(),
	)

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	log.Println("SMTP:", address)

	// Important: From must match SMTP Username
	return smtp.SendMail(address, auth, cfg.Username, []string{to}, msg)
}
