package appointment

import (
	"context"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models/appointment"
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/patient"
	"dental-clinic-system/models/user"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
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
	ParseTokenFromCookie(r *http.Request) (*claims.Claims, error)
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
func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedToken, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		helpers.WriteJSONError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	authenticatedUser, err := h.userService.GetUserByEmail(ctx, parsedToken.Email)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		helpers.WriteJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	appointments, err := h.appointmentService.GetAppointments(ctx, authenticatedUser.ClinicID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch appointments")
		helpers.WriteJSONError(w, "Failed to fetch appointments", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSONResponse(w, appointments, http.StatusOK)
}

// GetAppointment retrieves a single appointment by ID
func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Warn().Msgf("Invalid appointment ID: %s", idStr)
		helpers.WriteJSONError(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		helpers.WriteJSONError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	authenticatedUser, err := h.userService.GetUserByEmail(ctx, parsedToken.Email)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		helpers.WriteJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	retrievedAppointment, err := h.appointmentService.GetAppointment(ctx, uint(id))
	if err != nil {
		log.Error().Err(err).Msg("Appointment not found")
		helpers.WriteJSONError(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != retrievedAppointment.ClinicID {
		log.Warn().Msg("Unauthorized access to appointment")
		helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	helpers.WriteJSONResponse(w, retrievedAppointment, http.StatusOK)
}

// CreateAppointment creates a new appointment
func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Limit the size of the request body to prevent resource exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB

	var newAppointment appointment.Appointment
	if err := json.NewDecoder(r.Body).Decode(&newAppointment); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		helpers.WriteJSONError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		helpers.WriteJSONError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	authenticatedUser, err := h.userService.GetUserByEmail(ctx, parsedToken.Email)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		helpers.WriteJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	newAppointment.ClinicID = authenticatedUser.ClinicID

	createdAppointment, err := h.appointmentService.CreateAppointment(ctx, newAppointment)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create appointment")
		helpers.WriteJSONError(w, "Failed to create appointment", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSONResponse(w, createdAppointment, http.StatusCreated)
}

// UpdateAppointment updates an existing appointment
func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Limit the size of the request body
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB

	var updatedAppointment appointment.Appointment
	if err := json.NewDecoder(r.Body).Decode(&updatedAppointment); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		helpers.WriteJSONError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		helpers.WriteJSONError(w, "Invalid token", http.StatusUnauthorized)
		return
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
		helpers.WriteJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	if appointmentErr != nil {
		log.Error().Err(appointmentErr).Msg("Appointment not found")
		helpers.WriteJSONError(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != existingAppointment.ClinicID {
		log.Warn().Msg("Unauthorized access to update appointment")
		helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Ensure the appointment belongs to the clinic
	updatedAppointment.ClinicID = authenticatedUser.ClinicID

	updatedAppointment, err = h.appointmentService.UpdateAppointment(ctx, updatedAppointment)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update appointment")
		helpers.WriteJSONError(w, "Failed to update appointment", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSONResponse(w, updatedAppointment, http.StatusOK)
}

// DeleteAppointment deletes an appointment by ID
func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Warn().Msgf("Invalid appointment ID: %s", idStr)
		helpers.WriteJSONError(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		helpers.WriteJSONError(w, "Invalid token", http.StatusUnauthorized)
		return
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
		helpers.WriteJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	if appointmentErr != nil {
		log.Error().Err(appointmentErr).Msg("Appointment not found")
		helpers.WriteJSONError(w, "Appointment not found", http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != existingAppointment.ClinicID {
		log.Warn().Msg("Unauthorized access to delete appointment")
		helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.appointmentService.DeleteAppointment(ctx, uint(id)); err != nil {
		log.Error().Err(err).Msg("Failed to delete appointment")
		helpers.WriteJSONError(w, "Failed to delete appointment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// GetDoctorAppointments retrieves all appointments for a specific doctor
func (h *AppointmentHandler) GetDoctorAppointments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)
	idStr := params["id"]
	doctorID, err := strconv.Atoi(idStr)
	if err != nil || doctorID <= 0 {
		log.Warn().Msgf("Invalid doctor ID: %s", idStr)
		helpers.WriteJSONError(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		helpers.WriteJSONError(w, "Invalid token", http.StatusUnauthorized)
		return
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
		helpers.WriteJSONError(w, "Authenticated user not found", http.StatusNotFound)
		return
	}

	if doctorErr != nil {
		log.Error().Err(doctorErr).Msg("Doctor not found")
		helpers.WriteJSONError(w, "Doctor not found", http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != doctor.ClinicID {
		log.Warn().Msg("Forbidden access to doctor's appointments")
		helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	appointments, err := h.appointmentService.GetDoctorAppointments(ctx, uint(doctorID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch doctor's appointments")
		helpers.WriteJSONError(w, "Failed to fetch doctor's appointments", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSONResponse(w, appointments, http.StatusOK)
}

// GetPatientAppointments retrieves all appointments for a specific patient
func (h *AppointmentHandler) GetPatientAppointments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)
	idStr := params["id"]
	patientID, err := strconv.Atoi(idStr)
	if err != nil || patientID <= 0 {
		log.Warn().Msgf("Invalid patient ID: %s", idStr)
		helpers.WriteJSONError(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	parsedToken, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Error().Err(err).Msg("Invalid token")
		helpers.WriteJSONError(w, "Invalid token", http.StatusUnauthorized)
		return
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
		helpers.WriteJSONError(w, "Authenticated user not found", http.StatusNotFound)
		return
	}

	if patientErr != nil {
		log.Error().Err(patientErr).Msg("Patient not found")
		helpers.WriteJSONError(w, "Patient not found", http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != patientModel.ClinicID {
		log.Warn().Msg("Forbidden access to patient's appointments")
		helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}

	appointments, err := h.appointmentService.GetPatientAppointments(ctx, uint(patientID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch patient's appointments")
		helpers.WriteJSONError(w, "Failed to fetch patient's appointments", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSONResponse(w, appointments, http.StatusOK)
}
