package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AppointmentHandler struct {
	DB *gorm.DB
}

func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	var appointments []models.Appointment
	h.DB.Preload("Clinic").Preload("Patient").Preload("Doctor").Find(&appointments)
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
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		log.Printf("Error decoding appointment data: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()

	result := h.DB.Create(&appointment)
	if result.Error != nil {
		log.Printf("Error saving appointment to the database: %v", result.Error)
		http.Error(w, "Error saving appointment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
