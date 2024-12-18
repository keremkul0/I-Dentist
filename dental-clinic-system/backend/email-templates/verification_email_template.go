package email_templates

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
)

func CreateVerificationEmail(email, token string) (*gomail.Message, error) {

	//htmlBytes, err := os.ReadFile("./dental-clinic-system/Email_HTMLs/verification_email_html.html")
	//if err != nil {
	//	log.Fatalf("Şablon dosyası okunamadı: %v", err)
	//}
	//htmlTemplate := string(htmlBytes)

	tmpl, err := template.ParseFiles("..\\Email_HTMLs\\verification_email_html.html")
	if err != nil {
		log.Println("Şablon dosyası hatası:", err)
		return nil, err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, map[string]string{
		"VERIFY_LINK": "http://localhost:8080/verify-email?token=" + token,
	})
	if err != nil {
		log.Fatalf("Şablon işleme hatası: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@i-dentist.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "E-posta Doğrulama")
	m.SetBody("text/html", body.String())

	return m, nil
}
