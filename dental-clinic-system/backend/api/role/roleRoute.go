package role

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoleRoutes(router fiber.Router, handler *RoleHandler) {
	router.Get("/roles", handler.GetRoles)
	router.Get("/roles/:id", handler.GetRole)
}
