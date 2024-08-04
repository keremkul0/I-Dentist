package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "dental-clinic-system/models"
    "gorm.io/gorm"
)

type ClinicHandler struct {
    DB *gorm.DB
}

func (h *ClinicHandler) GetClinics(w http.ResponseWriter, r *http.Request) {
    var clinics []models.Clinic
    h.DB.Find(&clinics)
    json.NewEncoder(w).Encode(clinics)
}

func (h *ClinicHandler) GetClinic(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var clinic models.Clinic
    h.DB.First(&clinic, params["id"])
    json.NewEncoder(w).Encode(clinic)
}

func (h *ClinicHandler) CreateClinic(w http.ResponseWriter, r *http.Request) {
    var clinic models.Clinic
    json.NewDecoder(r.Body).Decode(&clinic)
    h.DB.Create(&clinic)
    json.NewEncoder(w).Encode(clinic)
}

func (h *ClinicHandler) UpdateClinic(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var clinic models.Clinic
    h.DB.First(&clinic, params["id"])
    json.NewDecoder(r.Body).Decode(&clinic)
    h.DB.Save(&clinic)
    json.NewEncoder(w).Encode(clinic)
}

func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var clinic models.Clinic
    h.DB.Delete(&clinic, params["id"])
    json.NewEncoder(w).Encode("Clinic deleted")
}
