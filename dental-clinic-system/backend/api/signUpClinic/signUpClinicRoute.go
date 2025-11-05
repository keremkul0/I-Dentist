package signUpClinic

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterSignupClinicRoutes(router fiber.Router, handler *SignUpClinicController) {
	router.Post("/singup-clinic", handler.SignUpClinic)
}
