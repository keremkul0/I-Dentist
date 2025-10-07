package resetPassword

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterResetPasswordRoutes(router fiber.Router, handler *ResetPasswordHandler) {
	router.Post("/reset-password", handler.ResetPassword)
}
