package procedureRepository

import (
	"context"
	"dental-clinic-system/models/procedure"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles procedure-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// GetProcedures retrieves all procedures for a specific clinic
func (repo *Repository) GetProcedures(ctx context.Context, clinicID uint) ([]procedure.Procedure, error) {
	var procs []procedure.Procedure
	result := repo.DB.WithContext(ctx).Where("clinic_id = ?", clinicID).Find(&procs)
	if result.Error != nil {
		log.Error().
			Str("operation", "GetProcedures").
			Err(result.Error).
			Uint("clinic_id", clinicID).
			Msg("Failed to retrieve procedures")
		return nil, result.Error
	}
	log.Info().
		Str("operation", "GetProcedures").
		Uint("clinic_id", clinicID).
		Int("count", len(procs)).
		Msg("Retrieved procedures successfully")
	return procs, nil
}

// GetProcedure retrieves a single procedure by its ID
func (repo *Repository) GetProcedure(ctx context.Context, id uint) (procedure.Procedure, error) {
	var proc procedure.Procedure
	result := repo.DB.WithContext(ctx).First(&proc, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Warn().
				Str("operation", "GetProcedure").
				Err(result.Error).
				Uint("procedure_id", id).
				Msg("Procedure not found")
		} else {
			log.Error().
				Str("operation", "GetProcedure").
				Err(result.Error).
				Uint("procedure_id", id).
				Msg("Failed to retrieve procedure")
		}
		return procedure.Procedure{}, result.Error
	}
	log.Info().
		Str("operation", "GetProcedure").
		Uint("procedure_id", id).
		Msg("Retrieved procedure successfully")
	return proc, nil
}

// CreateProcedure creates a new procedure record in the database
func (repo *Repository) CreateProcedure(ctx context.Context, newProc procedure.Procedure) (procedure.Procedure, error) {
	result := repo.DB.WithContext(ctx).Create(&newProc)
	if result.Error != nil {
		log.Error().
			Str("operation", "CreateProcedure").
			Err(result.Error).
			Msg("Failed to create procedure")
		return procedure.Procedure{}, result.Error
	}
	log.Info().
		Str("operation", "CreateProcedure").
		Uint("procedure_id", newProc.ID).
		Msg("Procedure created successfully")
	return newProc, nil
}

// UpdateProcedure updates an existing procedure record in the database
func (repo *Repository) UpdateProcedure(ctx context.Context, updatedProc procedure.Procedure) (procedure.Procedure, error) {
	result := repo.DB.WithContext(ctx).Save(&updatedProc)
	if result.Error != nil {
		log.Error().
			Str("operation", "UpdateProcedure").
			Err(result.Error).
			Uint("procedure_id", updatedProc.ID).
			Msg("Failed to update procedure")
		return procedure.Procedure{}, result.Error
	}
	log.Info().
		Str("operation", "UpdateProcedure").
		Uint("procedure_id", updatedProc.ID).
		Msg("Procedure updated successfully")
	return updatedProc, nil
}

// DeleteProcedure deletes a procedure record from the database by its ID
func (repo *Repository) DeleteProcedure(ctx context.Context, id uint) error {
	result := repo.DB.WithContext(ctx).Delete(&procedure.Procedure{}, id)
	if result.Error != nil {
		log.Error().
			Str("operation", "DeleteProcedure").
			Err(result.Error).
			Uint("procedure_id", id).
			Msg("Failed to delete procedure")
		return result.Error
	}
	log.Info().
		Str("operation", "DeleteProcedure").
		Uint("procedure_id", id).
		Msg("Procedure deleted successfully")
	return nil
}
