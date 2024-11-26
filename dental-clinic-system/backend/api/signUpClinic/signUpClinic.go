package signUpClinic

import (
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"
)

type SignUpClinicService interface {
	SignUpClinic(clinic models.Clinic, userCacheKey string) (models.Clinic, models.UserGetModel, error)
}

type SignUpClinicController struct {
	service SignUpClinicService
}

func NewSignUpClinicController(service SignUpClinicService) *SignUpClinicController {
	return &SignUpClinicController{service: service}
}

func (h *SignUpClinicController) SignUpClinic(w http.ResponseWriter, r *http.Request) {
	var clinic models.Clinic
	var id string

	err := json.NewDecoder(r.Body).Decode(&struct {
		Clinic *models.Clinic `json:"clinic"`
		Id     *string        `json:"id"`
	}{&clinic, &id})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clinic, userGetModel, err := h.service.SignUpClinic(clinic, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(struct {
		Clinic models.Clinic       `json:"clinic"`
		User   models.UserGetModel `json:"user"`
	}{clinic, userGetModel})

	if err != nil {
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}