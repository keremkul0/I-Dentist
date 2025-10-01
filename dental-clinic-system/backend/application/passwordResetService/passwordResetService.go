package passwordResetService

import (
	"context"
	"dental-clinic-system/models/token"
	"dental-clinic-system/models/user"
	"errors"
)

type EmailService interface {
	SendPasswordResetEmail(ctx context.Context, email, token string) error
}

type PasswordResetTokenRepository interface {
	CreatePasswordResetToken(ctx context.Context, email string) (token.PasswordResetToken, error)
	ValidatePasswordResetToken(ctx context.Context, tokenStr string, email string) (token.PasswordResetToken, error)
	MarkTokenAsUsed(ctx context.Context, tokenID uint) error
	InvalidateAllTokensForEmail(ctx context.Context, email string) error
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
	UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string) error
}

type PasswordResetService struct {
	emailService                 EmailService
	passwordResetTokenRepository PasswordResetTokenRepository
	userRepository               UserRepository
}

func NewPasswordResetService(emailService EmailService, passwordResetTokenRepository PasswordResetTokenRepository, userRepository UserRepository) *PasswordResetService {
	return &PasswordResetService{
		emailService:                 emailService,
		passwordResetTokenRepository: passwordResetTokenRepository,
		userRepository:               userRepository,
	}
}

// RequestPasswordReset - Şifre sıfırlama isteği oluşturma
func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, email string) error {
	// Kullanıcının varlığını kontrol et
	_, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	// Eski tokenları geçersiz kıl
	err = s.passwordResetTokenRepository.InvalidateAllTokensForEmail(ctx, email)
	if err != nil {
		return err
	}

	// Yeni token oluştur
	resetToken, err := s.passwordResetTokenRepository.CreatePasswordResetToken(ctx, email)
	if err != nil {
		return err
	}

	// Email gönder
	return s.emailService.SendPasswordResetEmail(ctx, email, resetToken.Token)
}

// ValidateResetToken - Token doğrulama
func (s *PasswordResetService) ValidateResetToken(ctx context.Context, tokenStr string, email string) (*token.PasswordResetToken, error) {
	resetToken, err := s.passwordResetTokenRepository.ValidatePasswordResetToken(ctx, tokenStr, email)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	return &resetToken, nil
}

// ResetPassword - Yeni şifre belirleme
func (s *PasswordResetService) ResetPassword(ctx context.Context, tokenStr string, email string, newHashedPassword string) error {
	// Token'ı doğrula
	resetToken, err := s.ValidateResetToken(ctx, tokenStr, email)
	if err != nil {
		return err
	}

	// Kullanıcıyı bul
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	// Şifreyi güncelle
	err = s.userRepository.UpdateUserPassword(ctx, user.ID, newHashedPassword)
	if err != nil {
		return err
	}

	// Token'ı kullanılmış olarak işaretle
	return s.passwordResetTokenRepository.MarkTokenAsUsed(ctx, resetToken.ID)
}
