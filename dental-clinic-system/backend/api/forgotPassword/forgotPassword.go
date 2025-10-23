package forgotPassword

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetService interface {
	RequestPasswordReset(ctx context.Context, email string) error
}

type ForgotPasswordHandler struct {
	passwordResetService PasswordResetService
}

func NewForgotPasswordController(passwordResetService PasswordResetService) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		passwordResetService: passwordResetService,
	}
}

func (h *ForgotPasswordHandler) ForgotPassword(c *fiber.Ctx) error {
	// Nil check for safety
	if c == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid context")
	}

	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	ctx := c.Context()
	if ctx == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid request context",
		})
	}

	// Error handling for password reset service
	err := h.passwordResetService.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		// Log error but don't expose details for security
		// log.Error().Err(err).Str("email", req.Email).Msg("Password reset request failed")
	}

	// Check if context is still valid before responding
	if c.Context().Err() != nil {
		return fiber.NewError(fiber.StatusGatewayTimeout, "Request timeout")
	}

	// Güvenlik nedeniyle her durumda aynı response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If the email exists, a password reset link has been sent",
	})
}
