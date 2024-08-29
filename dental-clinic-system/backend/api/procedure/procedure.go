package procedure

import (
	"dental-clinic-system/application/procedureService"
	"dental-clinic-system/models"
	"dental-clinic-system/repository/procedureRepository"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ProcedureHandlerService interface {
	GetProcedures(w http.ResponseWriter, r *http.Request)
	GetProcedure(w http.ResponseWriter, r *http.Request)
	CreateProcedure(w http.ResponseWriter, r *http.Request)
	UpdateProcedure(w http.ResponseWriter, r *http.Request)
	DeleteProcedure(w http.ResponseWriter, r *http.Request)
}

func NewProcedureHandlerService(db *gorm.DB) *ProcedureHandler {
	newProcedureRepository := procedureRepository.NewProcedureRepository(db)
	newProcedureService := procedureService.NewProcedureService(newProcedureRepository)
	return &ProcedureHandler{procedureService: newProcedureService}
}

type ProcedureHandler struct {
	procedureService procedureService.ProcedureService
}

func (h *ProcedureHandler) GetProcedures(w http.ResponseWriter, r *http.Request) {
	procedures, err := h.procedureService.GetProcedures()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(procedures)
	if err != nil {
		return
	}
}

func (h *ProcedureHandler) GetProcedure(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid procedure ID", http.StatusBadRequest)
		return
	}
	procedure, err := h.procedureService.GetProcedure(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(procedure)
	if err != nil {
		return
	}
}

func (h *ProcedureHandler) CreateProcedure(w http.ResponseWriter, r *http.Request) {
	var procedure models.Procedure
	err := json.NewDecoder(r.Body).Decode(&procedure)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	procedure, err = h.procedureService.CreateProcedure(procedure)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(procedure)
	if err != nil {
		return
	}
}

func (h *ProcedureHandler) UpdateProcedure(w http.ResponseWriter, r *http.Request) {
	var procedure models.Procedure
	err := json.NewDecoder(r.Body).Decode(&procedure)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	procedure, err = h.procedureService.UpdateProcedure(procedure)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(procedure)
	if err != nil {
		return
	}
}

func (h *ProcedureHandler) DeleteProcedure(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid procedure ID", http.StatusBadRequest)
		return
	}
	err = h.procedureService.DeleteProcedure(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
