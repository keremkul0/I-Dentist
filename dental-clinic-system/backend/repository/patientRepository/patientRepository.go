package patientRepository

import (
	"context"
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

func NewPatientRepository(db *gorm.DB) *patientRepository {
	return &patientRepository{DB: db}
}

type patientRepository struct {
	DB *gorm.DB
}

func (r *patientRepository) GetPatients(ctx context.Context, ClinicID uint) ([]models.Patient, error) {
	var patients []models.Patient
	result := r.DB.WithContext(ctx).Where("clinic_id = ?", ClinicID).Find(&patients)
	return patients, result.Error
}

func (r *patientRepository) GetPatient(ctx context.Context, id uint) (models.Patient, error) {
	var patient models.Patient
	result := r.DB.WithContext(ctx).First(&patient, id)
	return patient, result.Error
}

func (r *patientRepository) CreatePatient(ctx context.Context, patient models.Patient) (models.Patient, error) {
	result := r.DB.WithContext(ctx).Create(&patient)
	return patient, result.Error
}

func (r *patientRepository) UpdatePatient(ctx context.Context, patient models.Patient) (models.Patient, error) {
	result := r.DB.WithContext(ctx).Save(&patient)
	return patient, result.Error
}

func (r *patientRepository) DeletePatient(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.Patient{}, id)
	return result.Error
}
