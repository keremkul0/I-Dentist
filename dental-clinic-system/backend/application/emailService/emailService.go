package emailService

import (
	email_Templates "dental-clinic-system/email-Templates"
	"dental-clinic-system/models"
	"gopkg.in/gomail.v2"
	"time"
)

type UserRepository interface {
	GetUserByEmailRepo(email string) (models.User, error)
	UpdateUserRepo(user models.User) (models.User, error)
}

type TokenRepository interface {
	AddTokenToBlacklistRepo(token string, time time.Time)
	IsTokenBlacklistedRepo(token string) bool
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

func (s *emailService) SendVerificationEmail(email, token string) error {

	m, err := email_Templates.CreateVerificationEmail(email, token)

	if err != nil {
		return err
	}

	return s.Email.DialAndSend(&m)
}

func (s *emailService) VerifyUserEmail(token string, email string) bool {
	if s.tokenRepository.IsTokenBlacklistedRepo(token) {
		return false
	}

	user, err := s.userRepository.GetUserByEmailRepo(email)
	if err != nil {
		return false
	}

	user.EmailVerified = true
	_, err = s.userRepository.UpdateUserRepo(user)

	s.tokenRepository.AddTokenToBlacklistRepo(token, time.Now().Add(time.Minute*10))
	return true
}
