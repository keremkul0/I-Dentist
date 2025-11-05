package handler

import (
	"email-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

type EmailHandler struct {
	emailService service.EmailServiceInterface
}

func NewEmailHandler(emailService service.EmailServiceInterface) *EmailHandler {
	return &EmailHandler{emailService: emailService}
}

func (h *EmailHandler) SendEmail(c *fiber.Ctx) error {
	var req service.EmailMessage

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Type == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email type is required",
		})
	}

	err := h.emailService.SendEmail(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to send email",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Email sent successfully",
	})
}
