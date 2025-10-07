package singUpUser

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterSignupUserRoutes(router fiber.Router, handler *SignUpUserHandler) {

	router.Post("/signup-user", handler.SignUpUser)
}
