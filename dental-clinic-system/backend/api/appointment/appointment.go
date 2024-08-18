package appointment

import (
	appointmentService "dental-clinic-system/application/appointment"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"

	"net/http"
)

type AppointmentHandlerService interface {
	GetAppointments(w http.ResponseWriter, r *http.Request)
	GetAppointment(w http.ResponseWriter, r *http.Request)
	CreateAppointment(w http.ResponseWriter, r *http.Request)
	UpdateAppointment(w http.ResponseWriter, r *http.Request)
	DeleteAppointment(w http.ResponseWriter, r *http.Request)
}

func NewAppointmentHandlerService(appointmentService appointmentService.AppointmentService) *appointmentHandler {
	return &appointmentHandler{
		appointmentService: appointmentService,
	}
}

type appointmentHandler struct {
	appointmentService appointmentService.AppointmentService
}

func (h *appointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	appointments, err := h.appointmentService.GetAppointments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(appointments)
	if err != nil {
		return
	}
}

func (h *appointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	appointment, err := h.appointmentService.GetAppointment(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(appointment)
	if err != nil {
		return
	}
}

func (h *appointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	appointment, err := h.appointmentService.CreateAppointment(appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(appointment)
	if err != nil {
		return
	}
}

func (h *appointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	appointment, err := h.appointmentService.UpdateAppointment(appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(appointment)
	if err != nil {
		return
	}
}

func (h *appointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := h.appointmentService.DeleteAppointment(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
