package forgotPassword

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterForgotPasswordRoutes(router fiber.Router, handler *ForgotPasswordHandler) {
	router.Post("/forgot-password", handler.ForgotPassword)
}
