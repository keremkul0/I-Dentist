package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "gorm.io/gorm"
    "dental-clinic-system/models"
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
    json.NewEncoder(w).Encode(assistants)
}

func (h *AssistantHandler) GetAssistant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var assistant models.Assistant
    if result := h.DB.First(&assistant, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(assistant)
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
    json.NewEncoder(w).Encode(assistant)
}

func (h *AssistantHandler) UpdateAssistant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var assistant models.Assistant
    if result := h.DB.First(&assistant, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&assistant); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    h.DB.Save(&assistant)
    json.NewEncoder(w).Encode(assistant)
}

func (h *AssistantHandler) DeleteAssistant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    if result := h.DB.Delete(&models.Assistant{}, params["id"]); result.Error != nil {
        http.Error(w, result.Error.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
