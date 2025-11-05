package service

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type EmailMessage struct {
	To   string            `json:"to"`
	Type string            `json:"type"` // verification, password_reset, notification
	Data map[string]string `json:"data,omitempty"`
}

type Mailer interface {
	SendMail(message gomail.Message) error
}

// EmailServiceInterface defines the contract for email service
type EmailServiceInterface interface {
	SendEmail(msg EmailMessage) error
}

type EmailService struct {
	mailer Mailer
}

func NewEmailService(mailer Mailer) *EmailService {
	return &EmailService{
		mailer: mailer,
	}
}

func (s *EmailService) SendEmail(msg EmailMessage) error {
	if msg.Type == "verification" {
		return s.sendVerificationEmail(msg.To, msg.Data["token"])
	} else {
		return s.sendPasswordResetEmail(msg.To, msg.Data["token"])
	}
}

func (s *EmailService) sendVerificationEmail(email, token string) error {
	return s.sendTemplateEmail(email, "E-posta Doğrulama",
		"templates/verification_email.html", map[string]string{
			"VERIFY_LINK": os.Getenv("FRONTEND_URL") + "/verify-email?token=" + token,
		},
	)
}

func (s *EmailService) sendPasswordResetEmail(email, token string) error {
	return s.sendTemplateEmail(
		email,
		"Şifre Sıfırlama",
		"templates/password_reset_email.html",
		map[string]string{
			"RESET_LINK": os.Getenv("FRONTEND_URL") + "/reset-password?token=" + token,
		},
	)
}

//func (s *EmailService) sendNotificationEmail(to, subject, body string) error {
//	return s.sendPlainEmail(to, subject, body)
//}
//
//func (s *EmailService) sendPlainEmail(to, subject, body string) error {
//	m := gomail.NewMessage()
//	m.SetHeader("From", os.Getenv("SMTP_FROM"))
//	m.SetHeader("To", to)
//	m.SetHeader("Subject", subject)
//	m.SetBody("text/html", body)
//
//	return s.mailer.SendMail(*m)
//}

func (s *EmailService) sendTemplateEmail(to, subject, templateFile string, data map[string]string) error {
	// Template dosyasının tam yolunu oluştur
	templatePath := filepath.Join("templates", filepath.Base(templateFile))

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Template file processing error: %v", err)
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Printf("Template processing error: %v", err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	return s.mailer.SendMail(*m)
}
