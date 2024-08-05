package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AppointmentHandler struct {
	DB *gorm.DB
}

func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	var appointments []models.Appointment
	if result := h.DB.Find(&appointments); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(appointments)
}

func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var appointment models.Appointment
	if result := h.DB.First(&appointment, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(appointment)
}

func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result := h.DB.Create(&appointment); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(appointment)
}

func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var appointment models.Appointment
	if result := h.DB.First(&appointment, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.DB.Save(&appointment)
	json.NewEncoder(w).Encode(appointment)
}

func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if result := h.DB.Delete(&models.Appointment{}, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
