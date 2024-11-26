package appointment

import (
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"

	"net/http"
)

type UserService interface {
	GetUser(id uint) (models.UserGetModel, error)
	GetUserByEmail(email string) (models.UserGetModel, error)
}

type AppointmentService interface {
	GetAppointments(ClinicID uint) ([]models.Appointment, error)
	GetAppointment(id uint) (models.Appointment, error)
	CreateAppointment(appointment models.Appointment) (models.Appointment, error)
	UpdateAppointment(appointment models.Appointment) (models.Appointment, error)
	DeleteAppointment(id uint) error
	GetDoctorAppointments(id uint) ([]models.Appointment, error)
	GetPatientAppointments(id uint) ([]models.Appointment, error)
}

type PatientService interface {
	GetPatient(id uint) (models.Patient, error)
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

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	appointments, err := h.appointmentService.GetAppointments(user.ClinicID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(appointments)
	if err != nil {
		return
	}
}

func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointmentRepository ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	appointment, err := h.appointmentService.GetAppointment(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != appointment.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = json.NewEncoder(w).Encode(appointment)
	if err != nil {
		return
	}
}

func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	appointment.ClinicID = user.ClinicID

	appointment, err = h.appointmentService.CreateAppointment(appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(appointment)
	if err != nil {
		return
	}
}

func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	GetAppointment, err := h.appointmentService.GetAppointment(appointment.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != GetAppointment.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointment, err = h.appointmentService.UpdateAppointment(appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(appointment)
	if err != nil {
		return
	}
}

func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointmentRepository ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	GetAppointment, err := h.appointmentService.GetAppointment(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != GetAppointment.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.appointmentService.DeleteAppointment(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AppointmentHandler) GetDoctorAppointments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	doctor, err := h.userService.GetUser(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != doctor.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointments, err := h.appointmentService.GetDoctorAppointments(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(appointments)
	if err != nil {
		return
	}
}

func (h *AppointmentHandler) GetPatientAppointments(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	patient, err := h.patientService.GetPatient(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != patient.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointments, err := h.appointmentService.GetPatientAppointments(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(appointments)
	if err != nil {
		return
	}
}
