package clinic

import (
	"dental-clinic-system/middleware/rbacMiddleware"
	"dental-clinic-system/models/user"

	"github.com/gofiber/fiber/v2"
)

func RegisterClinicRoutes(router fiber.Router, handler *ClinicHandler) {
	//router.Get("/clinics", handler.GetClinics)
	router.Get("/clinic/:id", handler.GetClinic)
	//router.Post("/clinic", handler.CreateClinic)
	router.Put("/clinic", rbacMiddleware.RequireRole(user.RoleClinicAdmin, user.RoleSuperAdmin), handler.UpdateClinic)
	//router.Delete("/clinic/{id}", handler.DeleteClinic)
}
