package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "dental-clinic-system/models"
    "gorm.io/gorm"
)

type AppointmentHandler struct {
    DB *gorm.DB
}

func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
    var appointments []models.Appointment
    h.DB.Preload("Patient").Preload("Doctor").Preload("Clinic").Find(&appointments)
    json.NewEncoder(w).Encode(appointments)
}

func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var appointment models.Appointment
    h.DB.Preload("Patient").Preload("Doctor").Preload("Clinic").First(&appointment, params["id"])
    json.NewEncoder(w).Encode(appointment)
}

func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
    var appointment models.Appointment
    json.NewDecoder(r.Body).Decode(&appointment)
    h.DB.Create(&appointment)
    json.NewEncoder(w).Encode(appointment)
}

func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var appointment models.Appointment
    h.DB.First(&appointment, params["id"])
    json.NewDecoder(r.Body).Decode(&appointment)
    h.DB.Save(&appointment)
    json.NewEncoder(w).Encode(appointment)
}

func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var appointment models.Appointment
    h.DB.Delete(&appointment, params["id"])
    json.NewEncoder(w).Encode("Appointment deleted")
}
