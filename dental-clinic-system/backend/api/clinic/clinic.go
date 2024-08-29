package clinic

import (
	"dental-clinic-system/application/clinicService"
	"dental-clinic-system/models"
	"dental-clinic-system/repository/clinicRepository"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ClinicHandlerService interface {
	GetClinics(w http.ResponseWriter, r *http.Request)
	GetClinic(w http.ResponseWriter, r *http.Request)
	CreateClinic(w http.ResponseWriter, r *http.Request)
	UpdateClinic(w http.ResponseWriter, r *http.Request)
	GetClinicAppointments(w http.ResponseWriter, r *http.Request)
	DeleteClinic(w http.ResponseWriter, r *http.Request)
}

func NewClinicHandlerService(db *gorm.DB) *ClinicHandler {
	clinicRepository := clinicRepository.NewClinicRepository(db)
	newClinicService := clinicService.NewClinicService(clinicRepository)
	return &ClinicHandler{clinicService: newClinicService}
}

type ClinicHandler struct {
	clinicService clinicService.ClinicService
}

func (h *ClinicHandler) GetClinics(w http.ResponseWriter, r *http.Request) {
	clinics, err := h.clinicService.GetClinics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(clinics)
	if err != nil {
		return
	}
}

func (h *ClinicHandler) GetClinic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointmentRepository ID", http.StatusBadRequest)
		return
	}
	clinic, err := h.clinicService.GetClinic(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(clinic)
	if err != nil {
		return
	}
}

func (h *ClinicHandler) CreateClinic(w http.ResponseWriter, r *http.Request) {
	var clinic models.Clinic
	err := json.NewDecoder(r.Body).Decode(&clinic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	clinic, err = h.clinicService.CreateClinic(clinic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(clinic)
	if err != nil {
		return
	}
}

func (h *ClinicHandler) UpdateClinic(w http.ResponseWriter, r *http.Request) {
	var clinic models.Clinic
	err := json.NewDecoder(r.Body).Decode(&clinic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	clinic, err = h.clinicService.UpdateClinic(clinic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(clinic)
	if err != nil {
		return
	}
}

func (h *ClinicHandler) GetClinicAppointments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointmentRepository ID", http.StatusBadRequest)
		return
	}
	appointments, err := h.clinicService.GetClinicAppointments(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(appointments)
	if err != nil {
		return
	}
}
