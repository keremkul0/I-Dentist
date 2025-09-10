package clinic

import (
	"context"
	"dental-clinic-system/models/claims"
	"dental-clinic-system/models/clinic"
	"dental-clinic-system/models/user"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rs/zerolog/log"
)

// UserService interface
type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
}

type JwtService interface {
	ParseTokenFromCookie(r *http.Request) (*claims.Claims, error)
}

// ClinicService interface
type ClinicService interface {
	GetClinics(ctx context.Context) ([]clinic.Clinic, error)
	GetClinic(ctx context.Context, id uint) (clinic.Clinic, error)
	CreateClinic(ctx context.Context, clinic clinic.Clinic) (clinic.Clinic, error)
	UpdateClinic(ctx context.Context, clinic clinic.Clinic) (clinic.Clinic, error)
	DeleteClinic(ctx context.Context, id uint) error
	CheckClinicExist(ctx context.Context, cln clinic.Clinic) (bool, error)
}

type RoleService interface {
	UserHasRole(user user.UserGetModel, roleName string) bool
}

// ClinicHandler handles HTTP requests for clinics
type ClinicHandler struct {
	clinicService ClinicService
	userService   UserService
	roleService   RoleService
	jwtService    JwtService
}

// NewClinicHandlerController creates a new instance of ClinicHandler
func NewClinicHandlerController(clinicService ClinicService, userService UserService, roleService RoleService, jwtService JwtService) *ClinicHandler {
	return &ClinicHandler{clinicService: clinicService, userService: userService, roleService: roleService, jwtService: jwtService}
}

// GetClinics retrieves all clinics
func (h *ClinicHandler) GetClinics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	clinics, err := h.clinicService.GetClinics(ctx)
	if err != nil {
		log.Error().
			Str("operation", "GetClinics").
			Err(err).
			Msg("Failed to retrieve clinics")
		http.Error(w, "Failed to retrieve clinics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(clinics)
	if err != nil {
		log.Error().
			Str("operation", "GetClinics").
			Err(err).
			Msg("Failed to encode clinics to JSON")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetClinic retrieves a single clinic by its ID
func (h *ClinicHandler) GetClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().
			Str("operation", "GetClinic").
			Str("clinic_id", idStr).
			Msg("Invalid clinic ID")
		http.Error(w, "Invalid clinic ID", http.StatusBadRequest)
		return
	}

	// Extract authenticatedUser from cookie
	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Warn().
			Str("operation", "GetClinic").
			Err(err).
			Msg("Unauthorized access - invalid token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Warn().
			Str("operation", "GetClinic").
			Err(err).
			Msg("Unauthorized access - authenticatedUser not found")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if authenticatedUser.ClinicID != uint(id) && !h.roleService.UserHasRole(authenticatedUser, "Superadmin") {
		log.Warn().
			Str("operation", "GetClinic").
			Uint("user_clinic_id", authenticatedUser.ClinicID).
			Uint("requested_clinic_id", uint(id)).
			Msg("User is not authorized to access this clinic")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cln, err := h.clinicService.GetClinic(ctx, uint(id))
	if err != nil {
		if errors.Is(err, clinic.ErrClinicNotFound) {
			log.Warn().
				Str("operation", "GetClinic").
				Uint("clinic_id", uint(id)).
				Msg("Clinic not found")
			http.Error(w, "Clinic not found", http.StatusNotFound)
			return
		}
		log.Error().
			Str("operation", "GetClinic").
			Err(err).
			Uint("clinic_id", uint(id)).
			Msg("Failed to retrieve clinic")
		http.Error(w, "Failed to retrieve clinic", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cln)
	if err != nil {
		log.Error().
			Str("operation", "GetClinic").
			Err(err).
			Msg("Failed to encode clinic to JSON")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateClinic creates a new clinic after validation and existence check
func (h *ClinicHandler) CreateClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var cln clinic.Clinic
	err := json.NewDecoder(r.Body).Decode(&cln)
	if err != nil {
		log.Warn().
			Str("operation", "CreateClinic").
			Err(err).
			Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cln, err = h.clinicService.CreateClinic(ctx, cln)
	if err != nil {
		if errors.Is(err, clinic.ErrClinicAlreadyExists) {
			log.Warn().
				Str("operation", "CreateClinic").
				Str("clinic_email", cln.Email).
				Msg("Clinic already exists")
			http.Error(w, "Clinic already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, clinic.ErrClinicValidation) {
			log.Warn().
				Str("operation", "CreateClinic").
				Err(err).
				Msg("Clinic validation failed")
			http.Error(w, "Clinic validation failed", http.StatusBadRequest)
			return
		}
		log.Error().
			Str("operation", "CreateClinic").
			Err(err).
			Msg("Failed to create clinic")
		http.Error(w, "Failed to create clinic", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(cln)
	if err != nil {
		log.Error().
			Str("operation", "CreateClinic").
			Err(err).
			Msg("Failed to encode clinic to JSON")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateClinic updates an existing clinic after validation and existence check
func (h *ClinicHandler) UpdateClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var cln clinic.Clinic
	err := json.NewDecoder(r.Body).Decode(&cln)
	if err != nil {
		log.Warn().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract authenticatedUser from cookie
	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Warn().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Unauthorized access - invalid token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Warn().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Unauthorized access - authenticatedUser not found")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if cln.ID != authenticatedUser.ClinicID && !h.roleService.UserHasRole(authenticatedUser, "Superadmin") {
		log.Warn().
			Str("operation", "UpdateClinic").
			Uint("user_clinic_id", authenticatedUser.ClinicID).
			Uint("clinic_id", cln.ID).
			Msg("User is not authorized to update this clinic")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cln, err = h.clinicService.UpdateClinic(ctx, cln)
	if err != nil {
		if errors.Is(err, clinic.ErrClinicNotFound) {
			log.Warn().
				Str("operation", "UpdateClinic").
				Uint("clinic_id", cln.ID).
				Msg("Clinic not found")
			http.Error(w, "Clinic not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, clinic.ErrClinicValidation) {
			log.Warn().
				Str("operation", "UpdateClinic").
				Err(err).
				Msg("Clinic validation failed")
			http.Error(w, "Clinic validation failed", http.StatusBadRequest)
			return
		}
		log.Error().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Failed to update clinic")
		http.Error(w, "Failed to update clinic", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cln)
	if err != nil {
		log.Error().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Failed to encode clinic to JSON")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteClinic deletes a clinic by its ID after existence check
func (h *ClinicHandler) DeleteClinic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn().
			Str("operation", "DeleteClinic").
			Str("clinic_id", idStr).
			Msg("Invalid clinic ID")
		http.Error(w, "Invalid clinic ID", http.StatusBadRequest)
		return
	}

	// Extract authenticatedUser from cookie to authorize delete operation
	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Warn().
			Str("operation", "DeleteClinic").
			Err(err).
			Msg("Unauthorized access - invalid token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Warn().
			Str("operation", "DeleteClinic").
			Err(err).
			Msg("Unauthorized access - authenticatedUser not found")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if authenticatedUser.ClinicID != uint(id) && !h.roleService.UserHasRole(authenticatedUser, "Superadmin") {
		log.Warn().
			Str("operation", "DeleteClinic").
			Uint("user_clinic_id", authenticatedUser.ClinicID).
			Uint("clinic_id", uint(id)).
			Msg("User is not authorized to delete this clinic")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.clinicService.DeleteClinic(ctx, uint(id))
	if err != nil {
		if errors.Is(err, clinic.ErrClinicNotFound) {
			log.Warn().
				Str("operation", "DeleteClinic").
				Uint("clinic_id", uint(id)).
				Msg("Clinic not found")
			http.Error(w, "Clinic not found", http.StatusNotFound)
			return
		}
		log.Error().
			Str("operation", "DeleteClinic").
			Err(err).
			Uint("clinic_id", uint(id)).
			Msg("Failed to delete clinic")
		http.Error(w, "Failed to delete clinic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CheckClinicExistHandler checks if a clinic exists based on provided criteria
func (h *ClinicHandler) CheckClinicExist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var cln clinic.Clinic
	err := json.NewDecoder(r.Body).Decode(&cln)
	if err != nil {
		log.Warn().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract authenticatedUser from cookie
	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		log.Warn().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Unauthorized access - invalid token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		log.Warn().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Unauthorized access - authenticatedUser not found")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if cln.ID != authenticatedUser.ClinicID && !h.roleService.UserHasRole(authenticatedUser, "Superadmin") {
		log.Warn().
			Str("operation", "CheckClinicExist").
			Uint("user_clinic_id", authenticatedUser.ClinicID).
			Uint("clinic_id", cln.ID).
			Msg("User is not authorized to check this clinic")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exists, err := h.clinicService.CheckClinicExist(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Failed to check clinic existence")
		http.Error(w, "Failed to check clinic existence", http.StatusInternalServerError)
		return
	}

	response := map[string]bool{"exists": exists}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Failed to encode response to JSON")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
