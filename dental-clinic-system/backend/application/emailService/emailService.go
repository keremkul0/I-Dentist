package emailService

import (
	"context"
	email_Templates "dental-clinic-system/email-templates"
	"dental-clinic-system/models"
	"gopkg.in/gomail.v2"
	"time"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (models.User, error)
}

type TokenRepository interface {
	IsTokenBlacklisted(ctx context.Context, token string) bool
	AddTokenToBlacklist(ctx context.Context, token string, expireTime time.Time) error
}

type emailService struct {
	userRepository  UserRepository
	tokenRepository TokenRepository
	Email           gomail.Dialer
}

func NewEmailService(userRepository UserRepository, tokenRepository TokenRepository, Email gomail.Dialer) *emailService {
	return &emailService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		Email:           Email,
	}
}

func (s *emailService) SendVerificationEmail(ctx context.Context, email, token string) error {
	m, err := email_Templates.CreateVerificationEmail(email, token)
	if err != nil {
		return err
	}
	return s.Email.DialAndSend(m)
}

func (s *emailService) VerifyUserEmail(ctx context.Context, token string, email string) bool {
	if s.tokenRepository.IsTokenBlacklisted(ctx, token) {
		return false
	}

	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false
	}

	user.EmailVerified = true
	_, err = s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return false
	}

	err = s.tokenRepository.AddTokenToBlacklist(ctx, token, time.Now().Add(time.Minute*10))
	if err != nil {
		return false
	}
	return true
}
