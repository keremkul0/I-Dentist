package email

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	dialer *gomail.Dialer
}

func NewSMTPMailer() *SMTPMailer {
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	dialer := gomail.NewDialer(host, port, username, password)

	return &SMTPMailer{
		dialer: dialer,
	}
}

func (m *SMTPMailer) SendMail(message gomail.Message) error {
	return m.dialer.DialAndSend(&message)
}
