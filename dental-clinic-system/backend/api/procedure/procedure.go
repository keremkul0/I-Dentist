package procedure

import (
	"dental-clinic-system/application/procedureService"
	"dental-clinic-system/application/userService"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func NewProcedureController(procedureService procedureService.ProcedureService, userService userService.UserService) *ProcedureHandler {
	return &ProcedureHandler{
		procedureService: procedureService,
		userService:      userService,
	}
}

type ProcedureHandler struct {
	procedureService procedureService.ProcedureService
	userService      userService.UserService
}

func (h *ProcedureHandler) GetProcedures(w http.ResponseWriter, r *http.Request) {

	claims := helpers.TokenEmailHelper(r)
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	procedures, err := h.procedureService.GetProcedures(user.ClinicID)

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

	claims := helpers.TokenEmailHelper(r)
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	procedure, err := h.procedureService.GetProcedure(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if procedure.ClinicID != user.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	claims := helpers.TokenEmailHelper(r)
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	procedure.ClinicID = user.ClinicID

	procedure, err = h.procedureService.CreateProcedure(procedure)
	if err != nil {
		http.Error(w, "Failed to create procedure", http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(procedure)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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

	claims := helpers.TokenEmailHelper(r)
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !(helpers.ContainsRole(user, "Clinic Admin") || helpers.ContainsRole(user, "Superadmin")) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	procedure.ClinicID = user.ClinicID

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

	claims := helpers.TokenEmailHelper(r)
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !(helpers.ContainsRole(user, "Clinic Admin") || helpers.ContainsRole(user, "Superadmin")) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	procedure, err := h.procedureService.GetProcedure(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if procedure.ClinicID != user.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.procedureService.DeleteProcedure(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
