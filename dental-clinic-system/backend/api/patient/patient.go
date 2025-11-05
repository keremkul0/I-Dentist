package patient

import (
	"context"
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/patient"
	"dental-clinic-system/models/user"

	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
}

type PatientService interface {
	GetPatients(ctx context.Context, ClinicID uint) ([]patient.Patient, error)
	GetPatient(ctx context.Context, id uint) (patient.Patient, error)
	CreatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	UpdatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	DeletePatient(ctx context.Context, id uint) error
}

type JwtService interface {
	ParseTokenFromCookie(c *fiber.Ctx) (*claims.Claims, error)
}

type PatientHandler struct {
	patientService PatientService
	userService    UserService
	jwtService     JwtService
}

func NewPatientController(patientService PatientService, userService UserService, jwtService JwtService) *PatientHandler {
	return &PatientHandler{patientService: patientService, userService: userService, jwtService: jwtService}
}

func (h *PatientHandler) GetPatients(c *fiber.Ctx) error {
	ctx := c.Context()
	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	patients, err := h.patientService.GetPatients(ctx, user.ClinicID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(patients)
}

func (h *PatientHandler) GetPatient(c *fiber.Ctx) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid patient ID",
		})
	}

	patient, err := h.patientService.GetPatient(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if patient.ClinicID != user.ClinicID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(patient)
}

func (h *PatientHandler) CreatePatient(c *fiber.Ctx) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	var patient patient.Patient
	err := c.BodyParser(&patient)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	patient.ClinicID = user.ClinicID
	patient, err = h.patientService.CreatePatient(ctx, patient)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(patient)
}

func (h *PatientHandler) UpdatePatient(c *fiber.Ctx) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	var patient patient.Patient
	err := c.BodyParser(&patient)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	GetPatient, err := h.patientService.GetPatient(ctx, patient.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user.ClinicID != GetPatient.ClinicID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	patient, err = h.patientService.UpdatePatient(ctx, patient)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(patient)
}

func (h *PatientHandler) DeletePatient(c *fiber.Ctx) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid patient ID",
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	GetPatient, err := h.patientService.GetPatient(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user.ClinicID != GetPatient.ClinicID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	err = h.patientService.DeletePatient(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Patient deleted successfully",
	})
}
