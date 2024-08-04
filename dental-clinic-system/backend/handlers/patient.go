package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "dental-clinic-system/models"
    "gorm.io/gorm"
)

type PatientHandler struct {
    DB *gorm.DB
}

func (h *PatientHandler) GetPatients(w http.ResponseWriter, r *http.Request) {
    var patients []models.Patient
    h.DB.Find(&patients)
    json.NewEncoder(w).Encode(patients)
}

func (h *PatientHandler) GetPatient(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var patient models.Patient
    h.DB.First(&patient, params["id"])
    json.NewEncoder(w).Encode(patient)
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
    var patient models.Patient
    json.NewDecoder(r.Body).Decode(&patient)
    h.DB.Create(&patient)
    json.NewEncoder(w).Encode(patient)
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var patient models.Patient
    h.DB.First(&patient, params["id"])
    json.NewDecoder(r.Body).Decode(&patient)
    h.DB.Save(&patient)
    json.NewEncoder(w).Encode(patient)
}

func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var patient models.Patient
    h.DB.Delete(&patient, params["id"])
    json.NewEncoder(w).Encode("Patient deleted")
}
