package verifyEmail

import (
	"context"
	"dental-clinic-system/models/claims"

	"github.com/gofiber/fiber/v2"
)

type EmailService interface {
	VerifyUserEmail(ctx context.Context, token string, email string) bool
}

type JwtService interface {
	ParseToken(tokenStr string) (*claims.Claims, error)
}

type VerifyEmailHandler struct {
	EmailService EmailService
	jwtService   JwtService
}

func NewVerifyEmailController(service EmailService, jwtService JwtService) *VerifyEmailHandler {
	return &VerifyEmailHandler{EmailService: service, jwtService: jwtService}
}

func (h *VerifyEmailHandler) VerifyUserEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is required",
		})
	}

	claims, err := h.jwtService.ParseToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if h.EmailService.VerifyUserEmail(ctx, token, claims.Email) {
		return c.Status(fiber.StatusOK).SendString("Email verified")
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "Invalid token",
	})
}
