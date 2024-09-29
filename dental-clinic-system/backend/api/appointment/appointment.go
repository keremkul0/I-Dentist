package appointment

import (
	"dental-clinic-system/application/appointmentService"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"

	"net/http"
)

func NewAppointmentHandlerController(service appointmentService.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{appointmentService: service}
}

type AppointmentHandler struct {
	appointmentService appointmentService.AppointmentService
}

func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
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

func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointmentRepository ID", http.StatusBadRequest)
		return
	}
	appointment, err := h.appointmentService.GetAppointment(uint(id))
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

func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
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

func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointmentRepository ID", http.StatusBadRequest)
		return
	}
	err = h.appointmentService.DeleteAppointment(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
