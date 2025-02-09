package appointmentRepository

import (
	"context"
	"errors"

	"dental-clinic-system/models/appointment"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles appointment-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// GetAppointments retrieves all appointments for a specific clinic
func (repo *Repository) GetAppointments(ctx context.Context, clinicID uint) ([]appointment.Appointment, error) {
	var appointmentsList []appointment.Appointment
	result := repo.DB.WithContext(ctx).
		Where("clinic_id = ?", clinicID).
		Preload("Clinic").
		Preload("Patient").
		Preload("Doctor").
		Find(&appointmentsList)

	if result.Error != nil {
		log.Error().
			Str("operation", "GetAppointments").
			Err(result.Error).
			Uint("clinic_id", clinicID).
			Msg("Failed to retrieve appointments")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "GetAppointments").
		Uint("clinic_id", clinicID).
		Int("count", len(appointmentsList)).
		Msgf("Retrieved %d appointments successfully", len(appointmentsList))

	return appointmentsList, nil
}

// GetAppointment retrieves a single appointment by its ID
func (repo *Repository) GetAppointment(ctx context.Context, id uint) (appointment.Appointment, error) {
	var appt appointment.Appointment
	result := repo.DB.WithContext(ctx).
		Preload("Clinic").
		Preload("Patient").
		Preload("Doctor").
		First(&appt, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("operation", "GetAppointment").
				Err(result.Error).
				Uint("appointment_id", id).
				Msg("Appointment not found")
			return appointment.Appointment{}, gorm.ErrRecordNotFound
		}
		log.Error().
			Str("operation", "GetAppointment").
			Err(result.Error).
			Uint("appointment_id", id).
			Msg("Failed to retrieve appointment")
		return appointment.Appointment{}, result.Error
	}

	log.Info().
		Str("operation", "GetAppointment").
		Uint("appointment_id", id).
		Msg("Appointment retrieved successfully")

	return appt, nil
}

// CreateAppointment creates a new appointment record in the database
func (repo *Repository) CreateAppointment(ctx context.Context, newAppt appointment.Appointment) (appointment.Appointment, error) {
	result := repo.DB.WithContext(ctx).Create(&newAppt)
	if result.Error != nil {
		log.Error().
			Str("operation", "CreateAppointment").
			Err(result.Error).
			Msg("Failed to create appointment")
		return appointment.Appointment{}, result.Error
	}

	log.Info().
		Str("operation", "CreateAppointment").
		Uint("appointment_id", newAppt.ID).
		Msg("Appointment created successfully")

	return newAppt, nil
}

// UpdateAppointment updates an existing appointment record in the database
func (repo *Repository) UpdateAppointment(ctx context.Context, updatedAppt appointment.Appointment) (appointment.Appointment, error) {
	result := repo.DB.WithContext(ctx).Save(&updatedAppt)
	if result.Error != nil {
		log.Error().
			Str("operation", "UpdateAppointment").
			Err(result.Error).
			Uint("appointment_id", updatedAppt.ID).
			Msg("Failed to update appointment")
		return appointment.Appointment{}, result.Error
	}

	log.Info().
		Str("operation", "UpdateAppointment").
		Uint("appointment_id", updatedAppt.ID).
		Msg("Appointment updated successfully")

	return updatedAppt, nil
}

// DeleteAppointment deletes an appointment record from the database by its ID
func (repo *Repository) DeleteAppointment(ctx context.Context, id uint) error {
	result := repo.DB.WithContext(ctx).Delete(&appointment.Appointment{}, id)
	if result.Error != nil {
		log.Error().
			Str("operation", "DeleteAppointment").
			Err(result.Error).
			Uint("appointment_id", id).
			Msg("Failed to delete appointment")
		return result.Error
	}

	log.Info().
		Str("operation", "DeleteAppointment").
		Uint("appointment_id", id).
		Msg("Appointment deleted successfully")

	return nil
}

// GetDoctorAppointments retrieves all appointments for a specific doctor
func (repo *Repository) GetDoctorAppointments(ctx context.Context, doctorID uint) ([]appointment.Appointment, error) {
	var doctorAppointmentsList []appointment.Appointment
	result := repo.DB.WithContext(ctx).
		Where("doctor_id = ?", doctorID).
		Preload("Clinic").
		Preload("Patient").
		Preload("Doctor").
		Find(&doctorAppointmentsList)

	if result.Error != nil {
		log.Error().
			Str("operation", "GetDoctorAppointments").
			Err(result.Error).
			Uint("doctor_id", doctorID).
			Msg("Failed to retrieve doctor appointments")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "GetDoctorAppointments").
		Uint("doctor_id", doctorID).
		Int("count", len(doctorAppointmentsList)).
		Msgf("Retrieved %d appointments for doctor", len(doctorAppointmentsList))

	return doctorAppointmentsList, nil
}

// GetPatientAppointments retrieves all appointments for a specific patient
func (repo *Repository) GetPatientAppointments(ctx context.Context, patientID uint) ([]appointment.Appointment, error) {
	var patientAppointmentsList []appointment.Appointment
	result := repo.DB.WithContext(ctx).
		Where("patient_id = ?", patientID).
		Preload("Clinic").
		Preload("Patient").
		Preload("Doctor").
		Find(&patientAppointmentsList)

	if result.Error != nil {
		log.Error().
			Str("operation", "GetPatientAppointments").
			Err(result.Error).
			Uint("patient_id", patientID).
			Msg("Failed to retrieve patient appointments")
		return nil, result.Error
	}

	log.Info().
		Str("operation", "GetPatientAppointments").
		Uint("patient_id", patientID).
		Int("count", len(patientAppointmentsList)).
		Msgf("Retrieved %d appointments for patient", len(patientAppointmentsList))

	return patientAppointmentsList, nil
}
