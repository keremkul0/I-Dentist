package patientRepository

import (
	"context"
	"dental-clinic-system/models/patient"
	"errors"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles patient-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// GetPatients retrieves all patients for a specific clinic
func (repo *Repository) GetPatients(ctx context.Context, clinicID uint) ([]patient.Patient, error) {
	var patients []patient.Patient
	result := repo.DB.WithContext(ctx).Where("clinic_id = ?", clinicID).Find(&patients)
	if result.Error != nil {
		log.Error().
			Str("operation", "GetPatients").
			Err(result.Error).
			Uint("clinic_id", clinicID).
			Msg("Failed to retrieve patients")
		return nil, result.Error
	}
	log.Info().
		Str("operation", "GetPatients").
		Uint("clinic_id", clinicID).
		Int("count", len(patients)).
		Msg("Retrieved patients successfully")
	return patients, nil
}

// GetPatient retrieves a single patient by its ID
func (repo *Repository) GetPatient(ctx context.Context, id uint) (patient.Patient, error) {
	var pt patient.Patient
	result := repo.DB.WithContext(ctx).First(&pt, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("operation", "GetPatient").
				Err(result.Error).
				Uint("patient_id", id).
				Msg("Patient not found")
		} else {
			log.Error().
				Str("operation", "GetPatient").
				Err(result.Error).
				Uint("patient_id", id).
				Msg("Failed to retrieve patient")
		}
		return patient.Patient{}, result.Error
	}
	log.Info().
		Str("operation", "GetPatient").
		Uint("patient_id", id).
		Msg("Retrieved patient successfully")
	return pt, nil
}

// CreatePatient creates a new patient record in the database
func (repo *Repository) CreatePatient(ctx context.Context, newPt patient.Patient) (patient.Patient, error) {
	result := repo.DB.WithContext(ctx).Create(&newPt)
	if result.Error != nil {
		log.Error().
			Str("operation", "CreatePatient").
			Err(result.Error).
			Msg("Failed to create patient")
		return patient.Patient{}, result.Error
	}
	log.Info().
		Str("operation", "CreatePatient").
		Uint("patient_id", newPt.ID).
		Msg("Patient created successfully")
	return newPt, nil
}

// UpdatePatient updates an existing patient record in the database
func (repo *Repository) UpdatePatient(ctx context.Context, updatedPt patient.Patient) (patient.Patient, error) {
	result := repo.DB.WithContext(ctx).Save(&updatedPt)
	if result.Error != nil {
		log.Error().
			Str("operation", "UpdatePatient").
			Err(result.Error).
			Uint("patient_id", updatedPt.ID).
			Msg("Failed to update patient")
		return patient.Patient{}, result.Error
	}
	log.Info().
		Str("operation", "UpdatePatient").
		Uint("patient_id", updatedPt.ID).
		Msg("Patient updated successfully")
	return updatedPt, nil
}

// DeletePatient deletes a patient record from the database by its ID
func (repo *Repository) DeletePatient(ctx context.Context, id uint) error {
	result := repo.DB.WithContext(ctx).Delete(&patient.Patient{}, id)
	if result.Error != nil {
		log.Error().
			Str("operation", "DeletePatient").
			Err(result.Error).
			Uint("patient_id", id).
			Msg("Failed to delete patient")
		return result.Error
	}
	log.Info().
		Str("operation", "DeletePatient").
		Uint("patient_id", id).
		Msg("Patient deleted successfully")
	return nil
}
