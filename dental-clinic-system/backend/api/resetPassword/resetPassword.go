package resetPassword

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	Email       string `json:"email" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type PasswordResetService interface {
	ResetPassword(ctx context.Context, tokenStr string, email string, newHashedPassword string) error
}

type ResetPasswordHandler struct {
	passwordResetService PasswordResetService
}

func NewResetPasswordController(service PasswordResetService) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		passwordResetService: service,
	}
}

func (h *ResetPasswordHandler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest

	// JSON parse
	if err := c.BodyParser(&req); err != nil {
		log.Error().Err(err).Msg("Invalid JSON format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Validation
	if req.Token == "" || req.Email == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token, email and new_password are required",
		})
	}

	// Şifre hash'le
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Password processing failed",
		})
	}

	// Service çağrısı
	ctx := c.Context()
	err = h.passwordResetService.ResetPassword(ctx, req.Token, req.Email, string(hashedPassword))
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Password reset failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password reset failed",
		})
	}

	log.Info().Str("email", req.Email).Msg("Password reset successful")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successful",
	})
}
