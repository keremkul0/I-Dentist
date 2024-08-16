package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type AssistantHandler struct {
	DB *gorm.DB
}

func (h *AssistantHandler) GetAssistants(w http.ResponseWriter, r *http.Request) {
	var assistants []models.Assistant
	if result := h.DB.Find(&assistants); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(assistants)
	if err != nil {
		return
	}
}

func (h *AssistantHandler) GetAssistant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var assistant models.Assistant
	if result := h.DB.First(&assistant, params["id"]); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Assistant not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}
	err := json.NewEncoder(w).Encode(assistant)
	if err != nil {
		return
	}
}

func (h *AssistantHandler) CreateAssistant(w http.ResponseWriter, r *http.Request) {
	var assistant models.Assistant
	if err := json.NewDecoder(r.Body).Decode(&assistant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result := h.DB.Create(&assistant); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(assistant)
	if err != nil {
		return
	}
}

func (h *AssistantHandler) UpdateAssistant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var assistant models.Assistant
	if result := h.DB.First(&assistant, params["id"]); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Assistant not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&assistant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.DB.Save(&assistant)
	err := json.NewEncoder(w).Encode(assistant)
	if err != nil {
		return
	}
}

func (h *AssistantHandler) DeleteAssistant(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if result := h.DB.Delete(&models.Assistant{}, params["id"]); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Assistant not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
