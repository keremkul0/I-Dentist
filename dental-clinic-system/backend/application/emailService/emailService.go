package emailService

import (
	"context"
	"dental-clinic-system/models/user"
	"time"

	"github.com/rs/zerolog/log"
)

type EmailProducer interface {
	SendVerificationEmail(email string, token string) error
	SendPasswordResetEmail(email string, token string) error
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
	UpdateUser(ctx context.Context, user user.User) (user.User, error)
}

type TokenRepository interface {
	IsTokenBlacklisted(ctx context.Context, token string) bool
	AddTokenToBlacklist(ctx context.Context, token string, expireTime time.Time) error
}

type emailService struct {
	userRepository  UserRepository
	tokenRepository TokenRepository
	emailProducer   EmailProducer
}

func NewEmailService(userRepository UserRepository, tokenRepository TokenRepository, emailProducer EmailProducer) *emailService {
	return &emailService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		emailProducer:   emailProducer,
	}
}

func (s *emailService) SendVerificationEmail(email, token string) error {
	log.Info().Str("email", email).Msg("Sending verification email to Kafka")
	return s.emailProducer.SendVerificationEmail(email, token)
}

func (s *emailService) SendPasswordResetEmail(email, token string) error {
	log.Info().Str("email", email).Msg("Sending password reset email to Kafka")
	return s.emailProducer.SendPasswordResetEmail(email, token)
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
