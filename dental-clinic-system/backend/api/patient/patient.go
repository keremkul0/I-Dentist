package patient

import (
	"context"
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/patient"
	"dental-clinic-system/models/user"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
}

type PatientService interface {
	GetPatients(ctx context.Context, ClinicID uint) ([]patient.Patient, error)
	GetPatient(ctx context.Context, id uint) (patient.Patient, error)
	CreatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	UpdatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	DeletePatient(ctx context.Context, id uint) error
}

type JwtService interface {
	ParseTokenFromCookie(r *http.Request) (*claims.Claims, error)
}

type PatientHandler struct {
	patientService PatientService
	userService    UserService
	jwtService     JwtService
}

func NewPatientController(patientService PatientService, userService UserService, jwtService JwtService) *PatientHandler {
	return &PatientHandler{patientService: patientService, userService: userService, jwtService: jwtService}
}

func (h *PatientHandler) GetPatients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	patients, err := h.patientService.GetPatients(ctx, user.ClinicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(patients)
	if err != nil {
		return
	}
}

func (h *PatientHandler) GetPatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}
	patient, err := h.patientService.GetPatient(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if patient.ClinicID != user.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = json.NewEncoder(w).Encode(patient)
	if err != nil {
		return
	}
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	var patient patient.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	patient.ClinicID = user.ClinicID
	patient, err = h.patientService.CreatePatient(ctx, patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(patient)
	if err != nil {
		return
	}
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	var patient patient.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	GetPatient, err := h.patientService.GetPatient(ctx, patient.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != GetPatient.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	patient, err = h.patientService.UpdatePatient(ctx, patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(patient)
	if err != nil {
		return
	}
}

func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	GetPatient, err := h.patientService.GetPatient(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != GetPatient.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.patientService.DeletePatient(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
