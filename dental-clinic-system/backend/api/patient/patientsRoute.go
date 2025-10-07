package patient

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterPatientsRoutes(router fiber.Router, patientHandler *PatientHandler) {
	router.Get("/patients", patientHandler.GetPatients)
	router.Get("/patients/{id}", patientHandler.GetPatient)
	router.Post("/patients", patientHandler.CreatePatient)
	router.Put("/patients/{id}", patientHandler.UpdatePatient)
	router.Delete("/patients/{id}", patientHandler.DeletePatient)
}
