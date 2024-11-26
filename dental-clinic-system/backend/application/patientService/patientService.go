package patientService

import (
	"dental-clinic-system/models"
)

type PatientRepository interface {
	GetPatients(ClinicID uint) ([]models.Patient, error)
	GetPatient(id uint) (models.Patient, error)
	CreatePatient(patient models.Patient) (models.Patient, error)
	UpdatePatient(patient models.Patient) (models.Patient, error)
	DeletePatient(id uint) error
}

type patientService struct {
	patientRepository PatientRepository
}

func NewPatientService(patientRepository PatientRepository) *patientService {
	return &patientService{
		patientRepository: patientRepository,
	}
}

func (s *patientService) GetPatients(ClinicID uint) ([]models.Patient, error) {
	return s.patientRepository.GetPatients(ClinicID)
}

func (s *patientService) GetPatient(id uint) (models.Patient, error) {
	return s.patientRepository.GetPatient(id)
}

func (s *patientService) CreatePatient(patient models.Patient) (models.Patient, error) {
	return s.patientRepository.CreatePatient(patient)
}

func (s *patientService) UpdatePatient(patient models.Patient) (models.Patient, error) {
	return s.patientRepository.UpdatePatient(patient)
}

func (s *patientService) DeletePatient(id uint) error {
	return s.patientRepository.DeletePatient(id)
}
