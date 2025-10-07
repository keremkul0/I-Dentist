package procedure

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterProcedureRoutes(router fiber.Router, handler *ProcedureHandler) {
	router.Get("/procedures", handler.GetProcedures)
	router.Get("/procedures/{id}", handler.GetProcedure)
	router.Post("/procedures", handler.CreateProcedure)
	router.Put("/procedures/{id}", handler.UpdateProcedure)
	router.Delete("/procedures/{id}", handler.DeleteProcedure)
}
