package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type ProcedureHandler struct {
	DB *gorm.DB
}

func (h *ProcedureHandler) GetProcedures(w http.ResponseWriter, r *http.Request) {
	var procedures []models.Procedure
	if result := h.DB.Find(&procedures); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(procedures)
	if err != nil {
		return
	}
}

func (h *ProcedureHandler) GetProcedure(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var procedure models.Procedure
	if result := h.DB.First(&procedure, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(procedure)
	if err != nil {
		return
	}
}

func (h *ProcedureHandler) CreateProcedure(w http.ResponseWriter, r *http.Request) {
	var procedure models.Procedure
	if err := json.NewDecoder(r.Body).Decode(&procedure); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result := h.DB.Create(&procedure); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(procedure)
	if err != nil {
		return
	}
}

func (h *ProcedureHandler) UpdateProcedure(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var procedure models.Procedure
	if result := h.DB.First(&procedure, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&procedure); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.DB.Save(&procedure)
	err := json.NewEncoder(w).Encode(procedure)
	if err != nil {
		return
	}
}

func (h *ProcedureHandler) DeleteProcedure(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if result := h.DB.Delete(&models.Procedure{}, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
