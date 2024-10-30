package patient

import (
	"dental-clinic-system/application/patientService"
	"dental-clinic-system/application/userService"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func NewPatientController(patientService patientService.PatientService, userService userService.UserService) *PatientHandler {
	return &PatientHandler{patientService: patientService, userService: userService}
}

type PatientHandler struct {
	patientService patientService.PatientService
	userService    userService.UserService
}

func (h *PatientHandler) GetPatients(w http.ResponseWriter, r *http.Request) {
	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	patients, err := h.patientService.GetPatients(user.ClinicID)
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
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}
	patient, err := h.patientService.GetPatient(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

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
	var patient models.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	patient.ClinicID = user.ClinicID
	patient, err = h.patientService.CreatePatient(patient)
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
	var patient models.Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	GetPatient, err := h.patientService.GetPatient(patient.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != GetPatient.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	patient, err = h.patientService.UpdatePatient(patient)
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
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	GetPatient, err := h.patientService.GetPatient(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != GetPatient.ClinicID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.patientService.DeletePatient(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
