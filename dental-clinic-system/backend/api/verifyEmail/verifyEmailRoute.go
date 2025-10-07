package verifyEmail

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterVerifyEmailRoutes(router fiber.Router, handler *VerifyEmailHandler) {
	router.Get("/verify-email", handler.VerifyUserEmail)
}
