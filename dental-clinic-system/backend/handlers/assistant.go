package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "dental-clinic-system/models"
    "gorm.io/gorm"
)

type AssistantHandler struct {
    DB *gorm.DB
}

func (h *AssistantHandler) GetAssistants(w http.ResponseWriter, r *http.Request) {
    var assistants []models.Assistant
    h.DB.Preload("Clinic").Find(&assistants)
    json.NewEncoder(w).Encode(assistants)
}

func (h *AssistantHandler) GetAssistant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var assistant models.Assistant
    h.DB.Preload("Clinic").First(&assistant, params["id"])
    json.NewEncoder(w).Encode(assistant)
}

func (h *AssistantHandler) CreateAssistant(w http.ResponseWriter, r *http.Request) {
    var assistant models.Assistant
    json.NewDecoder(r.Body).Decode(&assistant)
    h.DB.Create(&assistant)
    json.NewEncoder(w).Encode(assistant)
}

func (h *AssistantHandler) UpdateAssistant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var assistant models.Assistant
    h.DB.First(&assistant, params["id"])
    json.NewDecoder(r.Body).Decode(&assistant)
    h.DB.Save(&assistant)
    json.NewEncoder(w).Encode(assistant)
}

func (h *AssistantHandler) DeleteAssistant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var assistant models.Assistant
    h.DB.Delete(&assistant, params["id"])
    json.NewEncoder(w).Encode("Assistant deleted")
}
