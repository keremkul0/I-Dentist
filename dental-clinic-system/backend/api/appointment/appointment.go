package appointment

import (
	"dental-clinic-system/repository/models"
	"encoding/json"
	"github.com/gorilla/mux"

	"net/http"

	appointmentService "dental-clinic-system/application/appointment"
)

type AppointmentHandlerService interface {
	GetAppointments(w http.ResponseWriter, r *http.Request)
	GetAppointment(w http.ResponseWriter, r *http.Request)
	CreateAppointment(w http.ResponseWriter, r *http.Request)
	UpdateAppointment(w http.ResponseWriter, r *http.Request)
	DeleteAppointment(w http.ResponseWriter, r *http.Request)
}

func NewAppointmentHandlerService(appointmentService *appointmentService.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		AppointmentService: appointmentService,
	}
}

type AppointmentHandler struct {
	AppointmentService *appointmentService.AppointmentService
}

func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	appointments, err := h.AppointmentService.GetAppointments()
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
	appointment, err := h.AppointmentService.GetAppointment(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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
	appointment, err := h.AppointmentService.CreateAppointment(appointment)
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
	appointment, err := h.AppointmentService.UpdateAppointment(appointment)
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
	err := h.AppointmentService.DeleteAppointment(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
