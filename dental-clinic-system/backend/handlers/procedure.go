package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "gorm.io/gorm"
    "dental-clinic-system/models"
)

type ProcedureHandler struct {
    DB *gorm.DB
}

func (h *ProcedureHandler) GetProcedures(w http.ResponseWriter, r *http.Request) {
    var procedures []models.Procedure
    h.DB.Find(&procedures)
    json.NewEncoder(w).Encode(&procedures)
}

func (h *ProcedureHandler) GetProcedure(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])
    var procedure models.Procedure
    if err := h.DB.First(&procedure, id).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(&procedure)
}

func (h *ProcedureHandler) CreateProcedure(w http.ResponseWriter, r *http.Request) {
    var procedure models.Procedure
    json.NewDecoder(r.Body).Decode(&procedure)
    h.DB.Create(&procedure)
    json.NewEncoder(w).Encode(&procedure)
}

func (h *ProcedureHandler) UpdateProcedure(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])
    var procedure models.Procedure
    if err := h.DB.First(&procedure, id).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    json.NewDecoder(r.Body).Decode(&procedure)
    h.DB.Save(&procedure)
    json.NewEncoder(w).Encode(&procedure)
}

func (h *ProcedureHandler) DeleteProcedure(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])
    var procedure models.Procedure
    if err := h.DB.First(&procedure, id).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    h.DB.Delete(&procedure)
    w.WriteHeader(http.StatusNoContent)
}
