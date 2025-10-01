package emailService

import (
	"bytes"
	"context"
	"dental-clinic-system/models/user"
	"github.com/rs/zerolog/log"
	"html/template"
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

type PasswordResetRepository interface {
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
	return s.sendTemplateEmail(ctx, email, "E-posta Doğrulama",
		"..\\Email_HTMLs\\verification_email_html.html", map[string]string{
			"VERIFY_LINK": "http://localhost:8080/verify-email?token=" + token,
		},
	)
}

func (s *emailService) SendPasswordResetEmail(ctx context.Context, email, token string) error {
	return s.sendTemplateEmail(
		ctx,
		email,
		"Şifre Sıfırlama",
		"..\\Email_HTMLs\\password_reset_email.html",
		map[string]string{
			"RESET_LINK": "http://localhost:8080/reset-password?token=" + token,
		},
	)
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

func (s *emailService) sendTemplateEmail(ctx context.Context, to, subject, templateFile string, data map[string]string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Error().Str("operation", "SendTemplateEmail").
			Err(err).
			Msg("Template file processing error")
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Error().Str("operation", "SendTemplateEmail").Err(err).Msg("Template proccessing error")
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@i-dentist.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	return s.mailer.SendMail(*m)
}
