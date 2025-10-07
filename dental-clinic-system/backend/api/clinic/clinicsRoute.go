package clinic

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterClinicRoutes(router fiber.Router, handler *ClinicHandler) {
	//router.Get("/clinics", handler.GetClinics)
	router.Get("/clinic/{id}", handler.GetClinic)
	//router.Post("/clinic", handler.CreateClinic)
	router.Put("/clinic/{id}", handler.UpdateClinic)
	//router.Delete("/clinic/{id}", handler.DeleteClinic)
}
