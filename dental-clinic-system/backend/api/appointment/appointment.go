package appointment

import (
	"context"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (models.UserGetModel, error)
}

type AppointmentService interface {
	GetAppointments(ctx context.Context, ClinicID uint) ([]models.Appointment, error)
	GetAppointment(ctx context.Context, id uint) (models.Appointment, error)
	CreateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error)
	UpdateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error)
	DeleteAppointment(ctx context.Context, id uint) error
	GetDoctorAppointments(ctx context.Context, id uint) ([]models.Appointment, error)
	GetPatientAppointments(ctx context.Context, id uint) ([]models.Appointment, error)
}

type PatientService interface {
	GetPatient(ctx context.Context, id uint) (models.Patient, error)
}

func NewAppointmentHandlerController(appointmentService AppointmentService, userService UserService, patientService PatientService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
		userService:        userService,
		patientService:     patientService,
	}
}

type AppointmentHandler struct {
	appointmentService AppointmentService
	userService        UserService
	patientService     PatientService
}

func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	appointments, err := h.appointmentService.GetAppointments(ctx, user.ClinicID)
	if err != nil {
		http.Error(w, "Appointments not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(appointments); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	appointment, err := h.appointmentService.GetAppointment(ctx, uint(id))
	if err != nil {
		http.Error(w, "Appointment not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != appointment.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(w).Encode(appointment); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	appointment.ClinicID = user.ClinicID

	appointment, err = h.appointmentService.CreateAppointment(ctx, appointment)
	if err != nil {
		http.Error(w, "Failed to create appointment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(appointment); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var user models.UserGetModel
	var getAppointment models.Appointment
	var userErr, appointmentErr error

	go func() {
		defer wg.Done()
		user, userErr = h.userService.GetUserByEmail(ctx, claims.Email)
	}()

	go func() {
		defer wg.Done()
		getAppointment, appointmentErr = h.appointmentService.GetAppointment(ctx, appointment.ID)
	}()

	wg.Wait()

	if userErr != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if appointmentErr != nil {
		http.Error(w, appointmentErr.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != getAppointment.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointment, err = h.appointmentService.UpdateAppointment(ctx, appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(appointment); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var user models.UserGetModel
	var getAppointment models.Appointment
	var userErr, appointmentErr error

	go func() {
		defer wg.Done()
		user, userErr = h.userService.GetUserByEmail(ctx, claims.Email)
	}()

	go func() {
		defer wg.Done()
		getAppointment, appointmentErr = h.appointmentService.GetAppointment(ctx, uint(id))
	}()

	wg.Wait()

	if userErr != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if appointmentErr != nil {
		http.Error(w, appointmentErr.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != getAppointment.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.appointmentService.DeleteAppointment(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AppointmentHandler) GetDoctorAppointments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var user models.UserGetModel
	var doctor models.UserGetModel
	var userErr, doctorErr error

	go func() {
		defer wg.Done()
		user, userErr = h.userService.GetUserByEmail(ctx, claims.Email)
	}()

	go func() {
		defer wg.Done()
		doctor, doctorErr = h.userService.GetUserByEmail(ctx, idStr)
	}()

	wg.Wait()

	if userErr != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if doctorErr != nil {
		http.Error(w, doctorErr.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != doctor.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointments, err := h.appointmentService.GetDoctorAppointments(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(appointments); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AppointmentHandler) GetPatientAppointments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var user models.UserGetModel
	var patient models.Patient
	var userErr, patientErr error

	go func() {
		defer wg.Done()
		user, userErr = h.userService.GetUserByEmail(ctx, claims.Email)
	}()

	go func() {
		defer wg.Done()
		patient, patientErr = h.patientService.GetPatient(ctx, uint(id))
	}()

	wg.Wait()

	if userErr != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if patientErr != nil {
		http.Error(w, patientErr.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != patient.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointments, err := h.appointmentService.GetPatientAppointments(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(appointments); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
