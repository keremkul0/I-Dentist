package patient

import (
	"dental-clinic-system/application/patientService"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type patientHandlerController interface {
	GetPatients(w http.ResponseWriter, r *http.Request)
	GetPatient(w http.ResponseWriter, r *http.Request)
	CreatePatient(w http.ResponseWriter, r *http.Request)
	UpdatePatient(w http.ResponseWriter, r *http.Request)
	DeletePatient(w http.ResponseWriter, r *http.Request)
}

func NewPatientController(patientService patientService.PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

type PatientHandler struct {
	patientService patientService.PatientService
}

func (h *PatientHandler) GetPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := h.patientService.GetPatients()
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
	err = h.patientService.DeletePatient(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
