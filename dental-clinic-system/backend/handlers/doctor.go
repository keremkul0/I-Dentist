package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "gorm.io/gorm"
    "dental-clinic-system/models"
)

type DoctorHandler struct {
    DB *gorm.DB
}

func (h *DoctorHandler) GetDoctors(w http.ResponseWriter, r *http.Request) {
    var doctors []models.Doctor
    if result := h.DB.Find(&doctors); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(doctors)
}

func (h *DoctorHandler) GetDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var doctor models.Doctor
    if result := h.DB.First(&doctor, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
    var doctor models.Doctor
    if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if result := h.DB.Create(&doctor); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var doctor models.Doctor
    if result := h.DB.First(&doctor, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    h.DB.Save(&doctor)
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    if result := h.DB.Delete(&models.Doctor{}, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
