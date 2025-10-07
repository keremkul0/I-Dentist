package user

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(router fiber.Router, handler *UserHandler) {
	router.Get("/users", handler.GetUsers)
	router.Get("/users/:id", handler.GetUser)
	router.Post("/users", handler.CreateUser)
	router.Put("/users/:id", handler.UpdateUser)
	router.Delete("/users/:id", handler.DeleteUser)
}
