package patientService

import (
	"context"
	"dental-clinic-system/models/patient"
)

type PatientRepository interface {
	GetPatients(ctx context.Context, ClinicID uint) ([]patient.Patient, error)
	GetPatient(ctx context.Context, id uint) (patient.Patient, error)
	CreatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	UpdatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	DeletePatient(ctx context.Context, id uint) error
}

type patientService struct {
	patientRepository PatientRepository
}

func NewPatientService(patientRepository PatientRepository) *patientService {
	return &patientService{
		patientRepository: patientRepository,
	}
}

func (s *patientService) GetPatients(ctx context.Context, ClinicID uint) ([]patient.Patient, error) {
	return s.patientRepository.GetPatients(ctx, ClinicID)
}

func (s *patientService) GetPatient(ctx context.Context, id uint) (patient.Patient, error) {
	return s.patientRepository.GetPatient(ctx, id)
}

func (s *patientService) CreatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error) {
	return s.patientRepository.CreatePatient(ctx, patient)
}

func (s *patientService) UpdatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error) {
	return s.patientRepository.UpdatePatient(ctx, patient)
}

func (s *patientService) DeletePatient(ctx context.Context, id uint) error {
	return s.patientRepository.DeletePatient(ctx, id)
}
