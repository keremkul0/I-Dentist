package clinicRepository

import (
	"context"
	"errors"

	"dental-clinic-system/models/clinic"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles clinic-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// GetClinics retrieves all clinics from the database
func (repo *Repository) GetClinics(ctx context.Context) ([]clinic.Clinic, error) {
	var clinicsList []clinic.Clinic
	result := repo.DB.WithContext(ctx).Find(&clinicsList)
	if result.Error != nil {
		log.Error().
			Str("operation", "GetClinics").
			Err(result.Error).
			Msg("Failed to retrieve clinics")
		return nil, result.Error
	}
	log.Info().
		Str("operation", "GetClinics").
		Int("count", len(clinicsList)).
		Msg("Retrieved clinics successfully")
	return clinicsList, nil
}

// GetClinic retrieves a single clinic by its ID
func (repo *Repository) GetClinic(ctx context.Context, id uint) (clinic.Clinic, error) {
	var cln clinic.Clinic
	result := repo.DB.WithContext(ctx).First(&cln, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("operation", "GetClinic").
				Err(result.Error).
				Uint("clinic_id", id).
				Msg("Clinic not found")
			return clinic.Clinic{}, clinic.ErrClinicNotFound
		}
		log.Error().
			Str("operation", "GetClinic").
			Err(result.Error).
			Uint("clinic_id", id).
			Msg("Failed to retrieve clinic")
		return clinic.Clinic{}, result.Error
	}
	log.Info().
		Str("operation", "GetClinic").
		Uint("clinic_id", id).
		Msg("Retrieved clinic successfully")
	return cln, nil
}

// CreateClinic creates a new clinic record in the database
func (repo *Repository) CreateClinic(ctx context.Context, newCln clinic.Clinic) (clinic.Clinic, error) {
	result := repo.DB.WithContext(ctx).Create(&newCln)
	if result.Error != nil {
		log.Error().
			Str("operation", "CreateClinic").
			Err(result.Error).
			Msg("Failed to create clinic")
		return clinic.Clinic{}, result.Error
	}
	log.Info().
		Str("operation", "CreateClinic").
		Uint("clinic_id", newCln.ID).
		Msg("Clinic created successfully")
	return newCln, nil
}

// UpdateClinic updates an existing clinic record in the database
func (repo *Repository) UpdateClinic(ctx context.Context, updatedCln clinic.Clinic) (clinic.Clinic, error) {
	result := repo.DB.WithContext(ctx).Save(&updatedCln)
	if result.Error != nil {
		log.Error().
			Str("operation", "UpdateClinic").
			Err(result.Error).
			Uint("clinic_id", updatedCln.ID).
			Msg("Failed to update clinic")
		return clinic.Clinic{}, result.Error
	}
	log.Info().
		Str("operation", "UpdateClinic").
		Uint("clinic_id", updatedCln.ID).
		Msg("Clinic updated successfully")
	return updatedCln, nil
}

// DeleteClinic deletes a clinic record from the database by its ID
func (repo *Repository) DeleteClinic(ctx context.Context, id uint) error {
	result := repo.DB.WithContext(ctx).Delete(&clinic.Clinic{}, id)
	if result.Error != nil {
		log.Error().
			Str("operation", "DeleteClinic").
			Err(result.Error).
			Uint("clinic_id", id).
			Msg("Failed to delete clinic")
		return result.Error
	}
	log.Info().
		Str("operation", "DeleteClinic").
		Uint("clinic_id", id).
		Msg("Clinic deleted successfully")
	return nil
}

// CheckClinicExist checks if a clinic exists based on ID, email, name, or phone number
func (repo *Repository) CheckClinicExist(ctx context.Context, clnModel clinic.Clinic) (bool, error) {
	var count int64
	result := repo.DB.WithContext(ctx).
		Model(&clinic.Clinic{}).
		Where("id = ? OR email = ? OR name = ? OR phone_number = ?", clnModel.ID, clnModel.Email, clnModel.Name, clnModel.PhoneNumber).
		Count(&count)

	if result.Error != nil {
		log.Error().
			Str("operation", "CheckClinicExist").
			Err(result.Error).
			Msg("Failed to check clinic existence")
		return false, clinic.ErrClinicExistenceCheck
	}

	exists := count > 0
	if exists {
		log.Info().
			Str("operation", "CheckClinicExist").
			Int64("count", count).
			Msgf("Clinic exists with ID: %d, Email: %s, Name: %s, Phone: %s", clnModel.ID, clnModel.Email, clnModel.Name, clnModel.PhoneNumber)
	} else {
		log.Info().
			Str("operation", "CheckClinicExist").
			Msg("Clinic does not exist")
	}

	return exists, nil
}
