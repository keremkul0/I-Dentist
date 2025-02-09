package clinicService

import (
	"context"
	"dental-clinic-system/models/clinic"
	"dental-clinic-system/validations"
	"errors"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// ClinicRepository defines the interface for clinic-related database operations
type ClinicRepository interface {
	GetClinics(ctx context.Context) ([]clinic.Clinic, error)
	GetClinic(ctx context.Context, id uint) (clinic.Clinic, error)
	CreateClinic(ctx context.Context, cln clinic.Clinic) (clinic.Clinic, error)
	UpdateClinic(ctx context.Context, cln clinic.Clinic) (clinic.Clinic, error)
	DeleteClinic(ctx context.Context, id uint) error
	CheckClinicExist(ctx context.Context, cln clinic.Clinic) (bool, error)
}

// ClinicService handles clinic-related business logic
type ClinicService struct {
	clinicRepository ClinicRepository
}

// NewClinicService creates a new instance of ClinicService
func NewClinicService(clinicRepo ClinicRepository) *ClinicService {
	return &ClinicService{
		clinicRepository: clinicRepo,
	}
}

// GetClinics retrieves all clinics
func (s *ClinicService) GetClinics(ctx context.Context) ([]clinic.Clinic, error) {
	log.Info().
		Str("operation", "GetClinics").
		Msg("Fetching all clinics")

	clinics, err := s.clinicRepository.GetClinics(ctx)
	if err != nil {
		log.Error().
			Str("operation", "GetClinics").
			Err(err).
			Msg("Failed to retrieve clinics")
		return nil, err
	}

	log.Info().
		Str("operation", "GetClinics").
		Int("count", len(clinics)).
		Msgf("Retrieved %d clinics successfully", len(clinics))

	return clinics, nil
}

// GetClinic retrieves a single clinic by its ID
func (s *ClinicService) GetClinic(ctx context.Context, id uint) (clinic.Clinic, error) {
	log.Info().
		Str("operation", "GetClinic").
		Uint("clinic_id", id).
		Msg("Fetching clinic by ID")

	cln, err := s.clinicRepository.GetClinic(ctx, id)
	if err != nil {
		if errors.Is(err, clinic.ErrClinicNotFound) {
			log.Warn().
				Str("operation", "GetClinic").
				Uint("clinic_id", id).
				Msg("Clinic not found")
			return clinic.Clinic{}, clinic.ErrClinicNotFound
		}
		log.Error().
			Str("operation", "GetClinic").
			Err(err).
			Uint("clinic_id", id).
			Msg("Failed to retrieve clinic")
		return clinic.Clinic{}, err
	}

	log.Info().
		Str("operation", "GetClinic").
		Uint("clinic_id", id).
		Msg("Clinic retrieved successfully")

	return cln, nil
}

// CreateClinic creates a new clinic after validation and existence check
func (s *ClinicService) CreateClinic(ctx context.Context, cln clinic.Clinic) (clinic.Clinic, error) {
	log.Info().
		Str("operation", "CreateClinic").
		Str("clinic_email", cln.Email).
		Msg("Starting clinic creation process")

	// Validate clinic data
	if err := validations.ClinicValidation(&cln); err != nil {
		log.Error().
			Str("operation", "CreateClinic").
			Err(err).
			Msg("Clinic validation failed")
		return clinic.Clinic{}, clinic.ErrClinicValidation
	}

	// Check if clinic already exists
	exists, err := s.clinicRepository.CheckClinicExist(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "CreateClinic").
			Err(err).
			Msg("Error while checking if clinic exists")
		return clinic.Clinic{}, clinic.ErrClinicExistenceCheck
	}
	if exists {
		log.Warn().
			Str("operation", "CreateClinic").
			Str("clinic_email", cln.Email).
			Msg("Clinic already exists")
		return clinic.Clinic{}, clinic.ErrClinicAlreadyExists
	}

	// Create clinic record in the database
	createdCln, err := s.clinicRepository.CreateClinic(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "CreateClinic").
			Err(err).
			Str("clinic_email", cln.Email).
			Msg("Failed to create clinic")
		return clinic.Clinic{}, clinic.ErrClinicCreation
	}

	log.Info().
		Str("operation", "CreateClinic").
		Uint("clinic_id", createdCln.ID).
		Msg("Clinic created successfully")

	return createdCln, nil
}

// UpdateClinic updates an existing clinic after validation and existence check
func (s *ClinicService) UpdateClinic(ctx context.Context, cln clinic.Clinic) (clinic.Clinic, error) {
	log.Info().
		Str("operation", "UpdateClinic").
		Uint("clinic_id", cln.ID).
		Msg("Starting clinic update process")

	// Validate clinic data
	if err := validations.ClinicValidation(&cln); err != nil {
		log.Error().
			Str("operation", "UpdateClinic").
			Err(err).
			Uint("clinic_id", cln.ID).
			Msg("Clinic validation failed")
		return clinic.Clinic{}, clinic.ErrClinicValidation
	}

	// Check if clinic exists
	exists, err := s.clinicRepository.CheckClinicExist(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "UpdateClinic").
			Err(err).
			Msg("Error while checking if clinic exists")
		return clinic.Clinic{}, clinic.ErrClinicExistenceCheck
	}
	if !exists {
		log.Warn().
			Str("operation", "UpdateClinic").
			Uint("clinic_id", cln.ID).
			Msg("Clinic does not exist")
		return clinic.Clinic{}, clinic.ErrClinicNotFound
	}

	// Update clinic record in the database
	updatedCln, err := s.clinicRepository.UpdateClinic(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "UpdateClinic").
			Err(err).
			Uint("clinic_id", cln.ID).
			Msg("Failed to update clinic")
		return clinic.Clinic{}, clinic.ErrClinicUpdate
	}

	log.Info().
		Str("operation", "UpdateClinic").
		Uint("clinic_id", updatedCln.ID).
		Msg("Clinic updated successfully")

	return updatedCln, nil
}

// DeleteClinic deletes a clinic by its ID after existence check
func (s *ClinicService) DeleteClinic(ctx context.Context, id uint) error {
	log.Info().
		Str("operation", "DeleteClinic").
		Uint("clinic_id", id).
		Msg("Starting clinic deletion process")

	// Check if clinic exists
	exists, err := s.clinicRepository.CheckClinicExist(ctx, clinic.Clinic{Model: gorm.Model{ID: id}})
	if err != nil {
		log.Error().
			Str("operation", "DeleteClinic").
			Err(err).
			Uint("clinic_id", id).
			Msg("Error while checking if clinic exists")
		return clinic.ErrClinicExistenceCheck
	}
	if !exists {
		log.Warn().
			Str("operation", "DeleteClinic").
			Uint("clinic_id", id).
			Msg("Clinic does not exist")
		return clinic.ErrClinicNotFound
	}

	// Delete clinic record from the database
	err = s.clinicRepository.DeleteClinic(ctx, id)
	if err != nil {
		log.Error().
			Str("operation", "DeleteClinic").
			Err(err).
			Uint("clinic_id", id).
			Msg("Failed to delete clinic")
		return clinic.ErrClinicDeletion
	}

	log.Info().
		Str("operation", "DeleteClinic").
		Uint("clinic_id", id).
		Msg("Clinic deleted successfully")

	return nil
}

// CheckClinicExist checks if a clinic exists based on provided criteria
func (s *ClinicService) CheckClinicExist(ctx context.Context, cln clinic.Clinic) (bool, error) {
	log.Info().
		Str("operation", "CheckClinicExist").
		Msg("Checking if clinic exists")

	exists, err := s.clinicRepository.CheckClinicExist(ctx, cln)
	if err != nil {
		log.Error().
			Str("operation", "CheckClinicExist").
			Err(err).
			Msg("Failed to check clinic existence")
		return false, err
	}

	if exists {
		log.Info().
			Str("operation", "CheckClinicExist").
			Msg("Clinic exists")
	} else {
		log.Info().
			Str("operation", "CheckClinicExist").
			Msg("Clinic does not exist")
	}

	return exists, nil
}
