package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
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
	err := json.NewEncoder(w).Encode(doctors)
	if err != nil {
		return
	}
}

func (h *DoctorHandler) GetDoctor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var doctor models.Doctor
	if result := h.DB.First(&doctor, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(doctor)
	if err != nil {
		return
	}
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
	err := json.NewEncoder(w).Encode(doctor)
	if err != nil {
		return
	}
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
	err := json.NewEncoder(w).Encode(doctor)
	if err != nil {
		return
	}
}

func (h *DoctorHandler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if result := h.DB.Delete(&models.Doctor{}, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
