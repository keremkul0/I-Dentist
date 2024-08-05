package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "gorm.io/gorm"
    "dental-clinic-system/models"
)

type ClinicHandler struct {
    DB *gorm.DB
}

func (h *ClinicHandler) GetClinics(w http.ResponseWriter, r *http.Request) {
    var clinics []models.Clinic
    if result := h.DB.Find(&clinics); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(clinics)
}

func (h *ClinicHandler) GetClinic(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var clinic models.Clinic
    if result := h.DB.First(&clinic, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(clinic)
}

func (h *ClinicHandler) CreateClinic(w http.ResponseWriter, r *http.Request) {
    var clinic models.Clinic
    if err := json.NewDecoder(r.Body).Decode(&clinic); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if result := h.DB.Create(&clinic); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(clinic)
}

func (h *ClinicHandler) UpdateClinic(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var clinic models.Clinic
    if result := h.DB.First(&clinic, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&clinic); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    h.DB.Save(&clinic)
    json.NewEncoder(w).Encode(clinic)
}

func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    if result := h.DB.Delete(&models.Clinic{}, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
