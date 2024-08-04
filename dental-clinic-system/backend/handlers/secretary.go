package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "dental-clinic-system/models"
    "gorm.io/gorm"
)

type SecretaryHandler struct {
    DB *gorm.DB
}

func (h *SecretaryHandler) GetSecretaries(w http.ResponseWriter, r *http.Request) {
    var secretaries []models.Secretary
    h.DB.Preload("Clinic").Find(&secretaries)
    json.NewEncoder(w).Encode(secretaries)
}

func (h *SecretaryHandler) GetSecretary(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var secretary models.Secretary
    h.DB.Preload("Clinic").First(&secretary, params["id"])
    json.NewEncoder(w).Encode(secretary)
}

func (h *SecretaryHandler) CreateSecretary(w http.ResponseWriter, r *http.Request) {
    var secretary models.Secretary
    json.NewDecoder(r.Body).Decode(&secretary)
    h.DB.Create(&secretary)
    json.NewEncoder(w).Encode(secretary)
}

func (h *SecretaryHandler) UpdateSecretary(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var secretary models.Secretary
    h.DB.First(&secretary, params["id"])
    json.NewDecoder(r.Body).Decode(&secretary)
    h.DB.Save(&secretary)
    json.NewEncoder(w).Encode(secretary)
}

func (h *SecretaryHandler) DeleteSecretary(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var secretary models.Secretary
    h.DB.Delete(&secretary, params["id"])
    json.NewEncoder(w).Encode("Secretary deleted")
}
