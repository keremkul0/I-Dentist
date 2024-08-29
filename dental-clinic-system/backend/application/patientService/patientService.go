package patientService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/patientRepository"
)

type PatientService interface {
	GetPatients() ([]models.Patient, error)
	GetPatient(id uint) (models.Patient, error)
	CreatePatient(patient models.Patient) (models.Patient, error)
	UpdatePatient(patient models.Patient) (models.Patient, error)
	DeletePatient(id uint) error
}

type patientService struct {
	patientRepository patientRepository.PatientRepository
}

func NewPatientService(patientRepository patientRepository.PatientRepository) *patientService {
	return &patientService{
		patientRepository: patientRepository,
	}
}

func (s *patientService) GetPatients() ([]models.Patient, error) {
	return s.patientRepository.GetPatients()
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
