package signUpClinic

import (
	"context"
	"dental-clinic-system/models/clinic"
	"dental-clinic-system/models/user"
	"encoding/json"
	"net/http"
)

// SignUpClinicService defines the interface for clinic signup service
type SignUpClinicService interface {
	SignUpClinic(ctx context.Context, clinic clinic.Clinic, userCacheKey string) (clinic.Clinic, user.UserGetModel, error)
}

// SignUpClinicController handles clinic signup HTTP requests
type SignUpClinicController struct {
	service SignUpClinicService
}

// NewSignUpClinicController creates a new instance of SignUpClinicController
func NewSignUpClinicController(service SignUpClinicService) *SignUpClinicController {
	return &SignUpClinicController{service: service}
}

// signUpClinicRequest represents the incoming JSON request structure
type signUpClinicRequest struct {
	Clinic *clinic.Clinic `json:"clinic"`
	ID     *string        `json:"id"`
}

// signUpClinicResponse represents the outgoing JSON response structure
type signUpClinicResponse struct {
	Clinic clinic.Clinic     `json:"clinic"`
	User   user.UserGetModel `json:"user"`
}

// SignUpClinic handles the HTTP request for signing up a clinic
func (h *SignUpClinicController) SignUpClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req signUpClinicRequest

	// Decode the incoming JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Clinic == nil || req.ID == nil {
		http.Error(w, "Clinic and ID information are required.", http.StatusBadRequest)
		return
	}

	// Call the service layer to process the signup
	updatedClinic, userGetModel, err := h.service.SignUpClinic(ctx, *req.Clinic, *req.ID)
	if err != nil {
		http.Error(w, "Clinic signup failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := signUpClinicResponse{
		Clinic: updatedClinic,
		User:   userGetModel,
	}

	// Set the Content-Type header and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Encode and send the response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to send response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
