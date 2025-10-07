package appointment

import (
	"context"
	"dental-clinic-system/models/appointment"
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/patient"
	"dental-clinic-system/models/user"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// UserService defines methods to interact with user data
type UserService interface {
	GetUser(ctx context.Context, id uint) (user.UserGetModel, error)
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
}

// AppointmentService defines methods to interact with appointment data
type AppointmentService interface {
	GetAppointments(ctx context.Context, ClinicID uint) ([]appointment.Appointment, error)
	GetAppointment(ctx context.Context, id uint) (appointment.Appointment, error)
	CreateAppointment(ctx context.Context, appointment appointment.Appointment) (appointment.Appointment, error)
	UpdateAppointment(ctx context.Context, appointment appointment.Appointment) (appointment.Appointment, error)
	DeleteAppointment(ctx context.Context, id uint) error
	GetDoctorAppointments(ctx context.Context, id uint) ([]appointment.Appointment, error)
	GetPatientAppointments(ctx context.Context, id uint) ([]appointment.Appointment, error)
}

type JwtService interface {
	ParseTokenFromCookie(c *fiber.Ctx) (*claims.Claims, error)
}

// PatientService defines methods to interact with patient data
type PatientService interface {
	GetPatient(ctx context.Context, id uint) (patient.Patient, error)
}

// AppointmentHandler handles appointment-related HTTP requests
type AppointmentHandler struct {
	appointmentService AppointmentService
	userService        UserService
	patientService     PatientService
	jwtService         JwtService
}

// NewAppointmentHandler creates a new AppointmentHandler
func NewAppointmentHandler(as AppointmentService, us UserService, ps PatientService, jwtService JwtService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: as,
		userService:        us,
		patientService:     ps,
		jwtService:         jwtService,
	}
}

// GetAppointments retrieves all appointments for the authenticated user's clinic
func (h *AppointmentHandler) GetAppointments(c *fiber.Ctx) error {
	ctx := c.Context()

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	appointments, err := h.appointmentService.GetAppointments(ctx, authenticatedUser.ClinicID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch appointments")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch appointments",
		})
	}

	return c.Status(fiber.StatusOK).JSON(appointments)
}

// GetAppointment retrieves a single appointment by ID
func (h *AppointmentHandler) GetAppointment(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Warn().Msgf("Invalid appointment ID: %s", idStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid appointment ID",
		})
	}

	claims, err := h.jwtService.ParseTokenFromCookie(c)
	parsedToken, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	retrievedAppointment, err := h.appointmentService.GetAppointment(ctx, uint(id))
	if err != nil {
		log.Error().Err(err).Msg("Appointment not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Appointment not found",
		})
	}

	if authenticatedUser.ClinicID != retrievedAppointment.ClinicID {
		log.Warn().Msg("Unauthorized access to appointment")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized access to appointment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(retrievedAppointment)
}

// CreateAppointment creates a new appointment
func (h *AppointmentHandler) CreateAppointment(c *fiber.Ctx) error {
	ctx := c.Context()

	var newAppointment appointment.Appointment
	if err := c.BodyParser(&newAppointment); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	authenticatedUser, err := h.userService.GetUserByEmail(ctx, parsedToken.Email)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	newAppointment.ClinicID = authenticatedUser.ClinicID

	createdAppointment, err := h.appointmentService.CreateAppointment(ctx, newAppointment)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create appointment")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create appointment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdAppointment)
}

// UpdateAppointment updates an existing appointment
func (h *AppointmentHandler) UpdateAppointment(c *fiber.Ctx) error {
	ctx := c.Context()

	var updatedAppointment appointment.Appointment
	if err := c.BodyParser(&updatedAppointment); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var authenticatedUser user.UserGetModel
	var existingAppointment appointment.Appointment
	var userErr, appointmentErr error

	go func() {
		defer wg.Done()
		authenticatedUser, userErr = h.userService.GetUserByEmail(ctx, parsedToken.Email)
	}()

	go func() {
		defer wg.Done()
		existingAppointment, appointmentErr = h.appointmentService.GetAppointment(ctx, updatedAppointment.ID)
	}()

	wg.Wait()

	if userErr != nil {
		log.Error().Err(userErr).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if appointmentErr != nil {
		log.Error().Err(appointmentErr).Msg("Appointment not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Appointment not found",
		})
	}

	if authenticatedUser.ClinicID != existingAppointment.ClinicID {
		log.Warn().Msg("Unauthorized access to update appointment")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	// Ensure the appointment belongs to the clinic
	updatedAppointment.ClinicID = authenticatedUser.ClinicID

	updatedAppointment, err = h.appointmentService.UpdateAppointment(ctx, updatedAppointment)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update appointment")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update appointment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedAppointment)
}

// DeleteAppointment deletes an appointment by ID
func (h *AppointmentHandler) DeleteAppointment(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Warn().Msgf("Invalid appointment ID: %s", idStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid appointment ID",
		})
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var authenticatedUser user.UserGetModel
	var existingAppointment appointment.Appointment
	var userErr, appointmentErr error

	go func() {
		defer wg.Done()
		authenticatedUser, userErr = h.userService.GetUserByEmail(ctx, parsedToken.Email)
	}()

	go func() {
		defer wg.Done()
		existingAppointment, appointmentErr = h.appointmentService.GetAppointment(ctx, uint(id))
	}()

	wg.Wait()

	if userErr != nil {
		log.Error().Err(userErr).Msg("User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if appointmentErr != nil {
		log.Error().Err(appointmentErr).Msg("Appointment not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Appointment not found",
		})
	}

	if authenticatedUser.ClinicID != existingAppointment.ClinicID {
		log.Warn().Msg("Unauthorized access to delete appointment")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	if err := h.appointmentService.DeleteAppointment(ctx, uint(id)); err != nil {
		log.Error().Err(err).Msg("Failed to delete appointment")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete appointment",
		})
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

// GetDoctorAppointments retrieves all appointments for a specific doctor
func (h *AppointmentHandler) GetDoctorAppointments(c *fiber.Ctx) error {
	ctx := c.Context()

	idStr := c.Params("id")
	doctorID, err := strconv.Atoi(idStr)
	if err != nil || doctorID <= 0 {
		log.Warn().Msgf("Invalid doctor ID: %s", idStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid doctor ID",
		})
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var authenticatedUser user.UserGetModel
	var doctor user.UserGetModel
	var userErr, doctorErr error

	go func() {
		defer wg.Done()
		authenticatedUser, userErr = h.userService.GetUserByEmail(ctx, parsedToken.Email)
	}()

	go func() {
		defer wg.Done()
		doctor, doctorErr = h.userService.GetUser(ctx, uint(doctorID))
	}()

	wg.Wait()

	if userErr != nil {
		log.Error().Err(userErr).Msg("Authenticated user not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Authenticated user not found",
		})
	}

	if doctorErr != nil {
		log.Error().Err(doctorErr).Msg("Doctor not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Doctor not found",
		})
	}

	if authenticatedUser.ClinicID != doctor.ClinicID {
		log.Warn().Msg("Forbidden access to doctor's appointments")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	appointments, err := h.appointmentService.GetDoctorAppointments(ctx, uint(doctorID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch doctor's appointments")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch doctor's appointments",
		})
	}

	return c.Status(fiber.StatusOK).JSON(appointments)
}

// GetPatientAppointments retrieves all appointments for a specific patient
func (h *AppointmentHandler) GetPatientAppointments(c *fiber.Ctx) error {
	ctx := c.Context()

	idStr := c.Params("id")
	patientID, err := strconv.Atoi(idStr)
	if err != nil || patientID <= 0 {
		log.Warn().Msgf("Invalid patient ID: %s", idStr)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid patient ID",
		})
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(c)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var authenticatedUser user.UserGetModel
	var patientModel patient.Patient
	var userErr, patientErr error

	go func() {
		defer wg.Done()
		authenticatedUser, userErr = h.userService.GetUserByEmail(ctx, parsedToken.Email)
	}()

	go func() {
		defer wg.Done()
		patientModel, patientErr = h.patientService.GetPatient(ctx, uint(patientID))
	}()

	wg.Wait()

	if userErr != nil {
		log.Error().Err(userErr).Msg("Authenticated user not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Authenticated user not found",
		})
	}

	if patientErr != nil {
		log.Error().Err(patientErr).Msg("Patient not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Patient not found",
		})
	}

	if authenticatedUser.ClinicID != patientModel.ClinicID {
		log.Warn().Msg("Forbidden access to patient's appointments")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}

	appointments, err := h.appointmentService.GetPatientAppointments(ctx, uint(patientID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch patient's appointments")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch patient's appointments",
		})
	}

	return c.Status(fiber.StatusOK).JSON(appointments)
}
