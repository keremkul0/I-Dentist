package logout

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterLogoutRoutes(router fiber.Router, handler *LogoutController) {
	router.Post("/logout", handler.Logout)
}
