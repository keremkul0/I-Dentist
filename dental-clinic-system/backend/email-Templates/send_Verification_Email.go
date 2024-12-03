package email_Templates

import (
	"gopkg.in/gomail.v2"
)

func CreateVerificationEmail(email, token string) (gomail.Message, error) {

	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@i-dentist.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "E-posta Doğrulama")
	m.SetBody("text/plain", "Doğrulama için aşağıdaki bağlantıya tıklayın:\n\n"+
		"http://localhost:8080/verify-email?token="+token)

	return *m, nil
}
