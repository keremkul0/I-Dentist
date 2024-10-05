package clinic

import (
	"dental-clinic-system/application/clinicService"
	"dental-clinic-system/application/userService"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func NewClinicHandlerController(clinicService clinicService.ClinicService, userService userService.UserService) *ClinicHandler {
	return &ClinicHandler{clinicService: clinicService, userService: userService}
}

type ClinicHandler struct {
	clinicService clinicService.ClinicService
	userService   userService.UserService
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
	claims := helpers.TokenEmailHelper(r)
	user, err := h.userService.GetUserByEmail(claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != uint(id) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	claims := helpers.TokenEmailHelper(r)
	user, err := h.userService.GetUserByEmail(claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if clinic.ID != user.ClinicID && helpers.ContainsRole(user, "Superadmin") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointmentRepository ID", http.StatusBadRequest)
		return
	}
	err = h.clinicService.DeleteClinic(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
