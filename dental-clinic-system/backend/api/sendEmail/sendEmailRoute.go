package sendEmail

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterSendEmailRoutes(router fiber.Router, handler *SendEmailHandler) {
	router.Post("/send-verification-email", handler.SendVerificationEmail)
}
