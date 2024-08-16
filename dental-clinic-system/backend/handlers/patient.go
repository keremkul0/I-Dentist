package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type PatientHandler struct {
	DB *gorm.DB
}

func (h *PatientHandler) GetPatients(w http.ResponseWriter, r *http.Request) {
	var patients []models.Patient
	h.DB.Find(&patients)
	err := json.NewEncoder(w).Encode(patients)
	if err != nil {
		return
	}
}

func (h *PatientHandler) GetPatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var patient models.Patient
	h.DB.First(&patient, params["id"])
	err := json.NewEncoder(w).Encode(patient)
	if err != nil {
		return
	}
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		return
	}
	h.DB.Create(&patient)
	err = json.NewEncoder(w).Encode(patient)
	if err != nil {
		return
	}
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var patient models.Patient
	h.DB.First(&patient, params["id"])
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		return
	}
	h.DB.Save(&patient)
	err = json.NewEncoder(w).Encode(patient)
	if err != nil {
		return
	}
}

func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var patient models.Patient
	h.DB.Delete(&patient, params["id"])
	err := json.NewEncoder(w).Encode("Patient deleted")
	if err != nil {
		return
	}
}
