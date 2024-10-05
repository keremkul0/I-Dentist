package patientRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type PatientRepository interface {
	GetPatients(ClinicID uint) ([]models.Patient, error)
	GetPatient(id uint) (models.Patient, error)
	CreatePatient(patient models.Patient) (models.Patient, error)
	UpdatePatient(patient models.Patient) (models.Patient, error)
	DeletePatient(id uint) error
}

func NewPatientRepository(db *gorm.DB) *patientRepository {
	return &patientRepository{DB: db}
}

type patientRepository struct {
	DB *gorm.DB
}

func (r *patientRepository) GetPatients(ClinicID uint) ([]models.Patient, error) {
	var patients []models.Patient
	err := r.DB.Where("clinic_id = ?", ClinicID).Find(&patients).Error
	return patients, err
}

func (r *patientRepository) GetPatient(id uint) (models.Patient, error) {
	var patient models.Patient
	err := r.DB.First(&patient, id).Error
	return patient, err
}

func (r *patientRepository) CreatePatient(patient models.Patient) (models.Patient, error) {
	err := r.DB.Create(&patient).Error
	return patient, err
}

func (r *patientRepository) UpdatePatient(patient models.Patient) (models.Patient, error) {
	err := r.DB.Save(&patient).Error
	return patient, err
}

func (r *patientRepository) DeletePatient(id uint) error {
	err := r.DB.Delete(&models.Patient{}, id).Error
	return err
}
