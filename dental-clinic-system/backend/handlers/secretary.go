package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type SecretaryHandler struct {
	DB *gorm.DB
}

func (h *SecretaryHandler) GetSecretaries(w http.ResponseWriter, r *http.Request) {
	var secretaries []models.Secretary
	if result := h.DB.Find(&secretaries); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(secretaries)
	if err != nil {
		return
	}
}

func (h *SecretaryHandler) GetSecretary(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var secretary models.Secretary
	if result := h.DB.First(&secretary, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(secretary)
	if err != nil {
		return
	}
}

func (h *SecretaryHandler) CreateSecretary(w http.ResponseWriter, r *http.Request) {
	var secretary models.Secretary
	if err := json.NewDecoder(r.Body).Decode(&secretary); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result := h.DB.Create(&secretary); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(secretary)
	if err != nil {
		return
	}
}

func (h *SecretaryHandler) UpdateSecretary(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var secretary models.Secretary
	if result := h.DB.First(&secretary, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&secretary); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.DB.Save(&secretary)
	err := json.NewEncoder(w).Encode(secretary)
	if err != nil {
		return
	}
}

func (h *SecretaryHandler) DeleteSecretary(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if result := h.DB.Delete(&models.Secretary{}, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
