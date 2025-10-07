package login

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(router fiber.Router, handler *LoginHandler) {
	router.Post("/login", handler.Login)
}
