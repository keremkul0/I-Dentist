package appointment

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterAppointmentRoutes(router fiber.Router, handler *AppointmentHandler) {
	router.Get("/appointments", handler.GetAppointments)
	router.Get("/appointments/:id", handler.GetAppointment)
	router.Post("/appointments", handler.CreateAppointment)
	router.Put("/appointment/:id", handler.UpdateAppointment)
	router.Delete("/appointment/:id", handler.DeleteAppointment)
}
