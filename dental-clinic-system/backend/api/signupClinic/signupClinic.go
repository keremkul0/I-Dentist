package signupClinic

import (
	"dental-clinic-system/application/signupClinicService"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"
)

type SignUpClinicHandler struct {
	service signupClinicService.SignUpClinicService
}

func NewSignUpClinicHandler(service signupClinicService.SignUpClinicService) *SignUpClinicHandler {
	return &SignUpClinicHandler{service: service}
}

func (h *SignUpClinicHandler) SignUpClinic(w http.ResponseWriter, r *http.Request) {
	var clinic models.Clinic
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&struct {
		Clinic *models.Clinic `json:"clinic"`
		User   *models.User   `json:"user"`
	}{&clinic, &user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = helpers.HashPassword(user.Password)

	clinic, userGetModel, err := h.service.SignUpClinic(clinic, user)
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
