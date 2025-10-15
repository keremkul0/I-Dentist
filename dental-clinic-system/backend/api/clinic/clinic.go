package clinic

import (
	"context"
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/clinic"
	"dental-clinic-system/models/user"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// UserService interface
type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
}

type JwtService interface {
	ParseTokenFromCookie(c *fiber.Ctx) (*claims.Claims, error)
}

// ClinicService interface
type ClinicService interface {
	GetClinics(ctx context.Context) ([]clinic.Clinic, error)
	GetClinic(ctx context.Context, id uint) (clinic.Clinic, error)
	CreateClinic(ctx context.Context, clinic clinic.Clinic) (clinic.Clinic, error)
	UpdateClinic(ctx context.Context, clinic clinic.Clinic) (clinic.Clinic, error)
	DeleteClinic(ctx context.Context, id uint) error
	CheckClinicExist(ctx context.Context, cln clinic.Clinic) (bool, error)
}

type RoleService interface {
	UserHasRole(user user.UserGetModel, roleName string) bool
}

// ClinicHandler handles HTTP requests for clinics
type ClinicHandler struct {
	clinicService ClinicService
	userService   UserService
	roleService   RoleService
	jwtService    JwtService
}

// NewClinicHandlerController creates a new instance of ClinicHandler
func NewClinicHandlerController(clinicService ClinicService, userService UserService, roleService RoleService, jwtService JwtService) *ClinicHandler {
	return &ClinicHandler{clinicService: clinicService, userService: userService, roleService: roleService, jwtService: jwtService}
}

// GetClinics retrieves all clinics
func (h *ClinicHandler) GetClinics(c *fiber.Ctx) error {
	ctx := c.Context()
	clinics, err := h.clinicService.GetClinics(ctx)
	if err != nil {
		log.Error().
			Str("operation", "GetClinics").
			Err(err).
			Msg("Failed to retrieve clinics")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve clinics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(clinics)
}

// GetClinic retrieves a single clinic by its ID
func (h *ClinicHandler) GetClinic(c *fiber.Ctx) error {
	ctx := c.Context()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().
			Str("operation", "GetClinic").
			Str("clinic_id", idStr).
			Msg("Invalid clinic ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid clinic ID",
		})
	}

	// Extract authenticatedUser from cookie
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Warn().
			Str("operation", "GetClinic").
			Msg("Unauthorized access - claims not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Warn().
			Str("operation", "GetClinic").
			Err(err).
			Msg("Unauthorized access - user not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if authenticatedUser.ClinicID != uint(id) && !h.roleService.UserHasRole(authenticatedUser, "Superadmin") {
		log.Warn().
			Str("operation", "GetClinic").
			Uint("user_clinic_id", authenticatedUser.ClinicID).
			Uint("requested_clinic_id", uint(id)).
			Msg("User is not authorized to access this clinic")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	cln, err := h.clinicService.GetClinic(ctx, uint(id))
	if err != nil {
		if errors.Is(err, clinic.ErrClinicNotFound) {
			log.Warn().
				Str("operation", "GetClinic").
				Uint("clinic_id", uint(id)).
				Msg("Clinic not found")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Clinic not found",
			})
		}
		log.Error().
			Str("operation", "GetClinic").
			Err(err).
			Uint("clinic_id", uint(id)).
			Msg("Failed to retrieve clinic")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve clinic",
		})
	}

	return c.Status(fiber.StatusOK).JSON(cln)
}

// CreateClinic creates a new clinic after validation and existence check
func (h *ClinicHandler) CreateClinic(c *fiber.Ctx) error {
	ctx := c.Context()

	var cln clinic.Clinic
	if err := c.BodyParser(&cln); err != nil {
		log.Warn().
			Str("operation", "CreateClinic").
			Err(err).
			Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	createdClinic, err := h.clinicService.CreateClinic(ctx, cln)
	if err != nil {
		if errors.Is(err, clinic.ErrClinicAlreadyExists) {
			log.Warn().
				Str("operation", "CreateClinic").
				Str("clinic_email", cln.Email).
				Msg("Clinic already exists")
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Clinic already exists",
			})
		}
		if errors.Is(err, clinic.ErrClinicValidation) {
			log.Warn().
				Str("operation", "CreateClinic").
				Err(err).
				Msg("Clinic validation failed")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Clinic validation failed",
			})
		}
		log.Error().
			Str("operation", "CreateClinic").
			Err(err).
			Msg("Failed to create clinic")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create clinic",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdClinic)
}

// UpdateClinic updates an existing clinic after validation and existence check
func (h *ClinicHandler) UpdateClinic(c *fiber.Ctx) error {
	ctx := c.Context()

	var cln clinic.Clinic
	if err := c.BodyParser(&cln); err != nil {
		log.Warn().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Extract authenticatedUser from cookie
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Warn().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Unauthorized access - invalid token")

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Warn().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Unauthorized access - user not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if cln.ID != authenticatedUser.ClinicID && !h.roleService.UserHasRole(authenticatedUser, "Superadmin") {
		log.Warn().
			Str("operation", "UpdateClinic").
			Uint("user_clinic_id", authenticatedUser.ClinicID).
			Uint("clinic_id", cln.ID).
			Msg("User is not authorized to update this clinic")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	updatedClinic, err := h.clinicService.UpdateClinic(ctx, cln)
	if err != nil {
		if errors.Is(err, clinic.ErrClinicNotFound) {
			log.Warn().
				Str("operation", "UpdateClinic").
				Uint("clinic_id", cln.ID).
				Msg("Clinic not found")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Clinic not found",
			})
		}
		if errors.Is(err, clinic.ErrClinicValidation) {
			log.Warn().
				Str("operation", "UpdateClinic").
				Err(err).
				Msg("Clinic validation failed")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Clinic validation failed",
			})
		}
		log.Error().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Failed to update clinic")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update clinic",
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedClinic)
}

// DeleteClinic deletes a clinic by its ID after existence check
func (h *ClinicHandler) DeleteClinic(c *fiber.Ctx) error {
	ctx := c.Context()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().
			Str("operation", "DeleteClinic").
			Str("clinic_id", idStr).
			Msg("Invalid clinic ID")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid clinic ID",
		})
	}

	// Extract authenticatedUser from cookie to authorize delete operation
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Warn().
			Str("operation", "DeleteClinic").
			Err(err).
			Msg("Unauthorized access - invalid token")

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Warn().
			Str("operation", "DeleteClinic").
			Err(err).
			Msg("Unauthorized access - user not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if authenticatedUser.ClinicID != uint(id) && !h.roleService.UserHasRole(authenticatedUser, "Superadmin") {
		log.Warn().
			Str("operation", "DeleteClinic").
			Uint("user_clinic_id", authenticatedUser.ClinicID).
			Uint("clinic_id", uint(id)).
			Msg("User is not authorized to delete this clinic")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	err = h.clinicService.DeleteClinic(ctx, uint(id))
	if err != nil {
		if errors.Is(err, clinic.ErrClinicNotFound) {
			log.Warn().
				Str("operation", "DeleteClinic").
				Uint("clinic_id", uint(id)).
				Msg("Clinic not found")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Clinic not found",
			})
		}
		log.Error().
			Str("operation", "DeleteClinic").
			Err(err).
			Uint("clinic_id", uint(id)).
			Msg("Failed to delete clinic")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete clinic",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// CheckClinicExistHandler checks if a clinic exists based on provided criteria
func (h *ClinicHandler) CheckClinicExist(c *fiber.Ctx) error {
	ctx := c.Context()

	var cln clinic.Clinic
	if err := c.BodyParser(&cln); err != nil {
		log.Warn().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Extract authenticatedUser from cookie
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Warn().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Unauthorized access - invalid token")

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Warn().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Unauthorized access - user not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if cln.ID != authenticatedUser.ClinicID && !h.roleService.UserHasRole(authenticatedUser, "Superadmin") {
		log.Warn().
			Str("operation", "CheckClinicExist").
			Uint("user_clinic_id", authenticatedUser.ClinicID).
			Uint("clinic_id", cln.ID).
			Msg("User is not authorized to check this clinic")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	exists, err := h.clinicService.CheckClinicExist(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Failed to check clinic existence")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check clinic existence",
		})
	}

	response := map[string]bool{"exists": exists}
	return c.Status(fiber.StatusOK).JSON(response)
}
