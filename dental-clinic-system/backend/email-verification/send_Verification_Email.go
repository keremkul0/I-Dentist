package email_verification

import (
	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@i-dentist.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "E-posta Doğrulama")
	m.SetBody("text/plain", "Doğrulama için aşağıdaki bağlantıya tıklayın:\n\n"+
		"http://localhost:3000/verify-email?token="+token)

	d := gomail.NewDialer("smtp.example.com", 587, "username", "password")

	return d.DialAndSend(m)
}
