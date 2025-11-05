package sendEmail

import (
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/user"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EmailService interface {
	SendVerificationEmail(email string, token string) error
}

type JwtService interface {
	GenerateJWTToken(email string, roles []*user.Role, time time.Time) (string, error)
	ParseTokenFromCookie(c *fiber.Ctx) (*claims.Claims, error)
}

type SendEmailHandler struct {
	EmailService EmailService
	jwtService   JwtService
}

func NewSendEmailController(service EmailService, jwtService JwtService) *SendEmailHandler {
	return &SendEmailHandler{EmailService: service, jwtService: jwtService}
}

func (h *SendEmailHandler) SendVerificationEmail(c *fiber.Ctx) error {
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := h.jwtService.GenerateJWTToken(claims.Email, claims.Roles, time.Now().Add(time.Minute*5))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.EmailService.SendVerificationEmail(claims.Email, token)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Verification email sent successfully",
	})
}
