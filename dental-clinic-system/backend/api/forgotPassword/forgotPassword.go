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
	_ = h.passwordResetService.RequestPasswordReset(ctx, req.Email)

	// Güvenlik nedeniyle her durumda aynı response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "If the email exists, a password reset link has been sent",
	})
}
