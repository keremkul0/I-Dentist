package emailService

import (
	"bytes"
	"context"
	"dental-clinic-system/models/user"
	"html/template"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
	UpdateUser(ctx context.Context, user user.User) (user.User, error)
}

type TokenRepository interface {
	IsTokenBlacklisted(ctx context.Context, token string) bool
	AddTokenToBlacklist(ctx context.Context, token string, expireTime time.Time) error
}

type Mailer interface {
	SendMail(massage gomail.Message) error
}

type emailService struct {
	userRepository  UserRepository
	tokenRepository TokenRepository
	mailer          Mailer
}

func NewEmailService(userRepository UserRepository, tokenRepository TokenRepository, mailer Mailer) *emailService {
	return &emailService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		mailer:          mailer,
	}
}

func (s *emailService) SendVerificationEmail(ctx context.Context, email, token string) error {
	m, err := s.createVerificationEmail(email, token)
	if err != nil {
		return err
	}
	return s.mailer.SendMail(*m)
}

func (s *emailService) VerifyUserEmail(ctx context.Context, token string, email string) bool {
	if s.tokenRepository.IsTokenBlacklisted(ctx, token) {
		return false
	}

	foundUser, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false
	}

	foundUser.EmailVerified = true
	_, err = s.userRepository.UpdateUser(ctx, foundUser)
	if err != nil {
		return false
	}

	err = s.tokenRepository.AddTokenToBlacklist(ctx, token, time.Now().Add(time.Minute*10))
	if err != nil {
		return false
	}
	return true
}

func (s *emailService) createVerificationEmail(email, token string) (*gomail.Message, error) {

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
