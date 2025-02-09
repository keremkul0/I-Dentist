package procedure

import (
	"context"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models/procedure"
	"dental-clinic-system/models/user"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type ProcedureService interface {
	GetProcedures(ctx context.Context, ClinicID uint) ([]procedure.Procedure, error)
	GetProcedure(ctx context.Context, id uint) (procedure.Procedure, error)
	CreateProcedure(ctx context.Context, procedure procedure.Procedure) (procedure.Procedure, error)
	UpdateProcedure(ctx context.Context, procedure procedure.Procedure) (procedure.Procedure, error)
	DeleteProcedure(ctx context.Context, id uint) error
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
}

func NewProcedureController(procedureService ProcedureService, userService UserService) *ProcedureHandler {
	return &ProcedureHandler{
		procedureService: procedureService,
		userService:      userService,
	}
}

type ProcedureHandler struct {
	procedureService ProcedureService
	userService      UserService
}

func (h *ProcedureHandler) GetProcedures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	procedures, err := h.procedureService.GetProcedures(ctx, user.ClinicID)
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
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid procedure ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	procedure, err := h.procedureService.GetProcedure(ctx, uint(id))
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
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	var procedure procedure.Procedure
	err := json.NewDecoder(r.Body).Decode(&procedure)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	procedure.ClinicID = user.ClinicID
	procedure, err = h.procedureService.CreateProcedure(ctx, procedure)
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
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	var procedure procedure.Procedure
	err := json.NewDecoder(r.Body).Decode(&procedure)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !(helpers.ContainsRole(user, "Clinic Admin") || helpers.ContainsRole(user, "Superadmin")) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	procedure.ClinicID = user.ClinicID
	procedure, err = h.procedureService.UpdateProcedure(ctx, procedure)
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
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid procedure ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !(helpers.ContainsRole(user, "Clinic Admin") || helpers.ContainsRole(user, "Superadmin")) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	procedure, err := h.procedureService.GetProcedure(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if procedure.ClinicID != user.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.procedureService.DeleteProcedure(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
