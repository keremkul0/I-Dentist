package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ClinicHandler struct {
	DB *gorm.DB
}

func (h *ClinicHandler) CreateClinic(w http.ResponseWriter, r *http.Request) {
	var clinic models.Clinic
	err := json.NewDecoder(r.Body).Decode(&clinic)
	if err != nil {
		log.Printf("Error decoding clinic data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := h.DB.Create(&clinic)
	if result.Error != nil {
		log.Printf("Error saving clinic to the database: %v", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Clinic %s created successfully!", clinic.Name)
	w.WriteHeader(http.StatusCreated) // Status code changed to 201 Created
	err = json.NewEncoder(w).Encode(clinic)
	if err != nil {
		return
	}
}

func (h *ClinicHandler) GetClinics(w http.ResponseWriter, r *http.Request) {
	var clinics []models.Clinic
	if result := h.DB.Find(&clinics); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(clinics)
	if err != nil {
		return
	}
}

func (h *ClinicHandler) GetClinic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var clinic models.Clinic
	if result := h.DB.First(&clinic, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(clinic)
	if err != nil {
		return
	}
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
	err := json.NewEncoder(w).Encode(clinic)
	if err != nil {
		return
	}
}

func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if result := h.DB.Delete(&models.Clinic{}, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
