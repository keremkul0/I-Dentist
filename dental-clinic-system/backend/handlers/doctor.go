package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "dental-clinic-system/models"
    "gorm.io/gorm"
)

type DoctorHandler struct {
    DB *gorm.DB
}

func (h *DoctorHandler) GetDoctors(w http.ResponseWriter, r *http.Request) {
    var doctors []models.Doctor
    h.DB.Find(&doctors)
    json.NewEncoder(w).Encode(doctors)
}

func (h *DoctorHandler) GetDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var doctor models.Doctor
    h.DB.First(&doctor, params["id"])
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
    var doctor models.Doctor
    json.NewDecoder(r.Body).Decode(&doctor)
    h.DB.Create(&doctor)
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var doctor models.Doctor
    h.DB.First(&doctor, params["id"])
    json.NewDecoder(r.Body).Decode(&doctor)
    h.DB.Save(&doctor)
    json.NewEncoder(w).Encode(doctor)
}

func (h *DoctorHandler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var doctor models.Doctor
    h.DB.Delete(&doctor, params["id"])
    json.NewEncoder(w).Encode("Doctor deleted")
}
