package clinicService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/clinicRepository"
)

type ClinicService interface {
	GetClinics() ([]models.Clinic, error)
	GetClinic(id uint) (models.Clinic, error)
	CreateClinic(clinic models.Clinic) (models.Clinic, error)
	UpdateClinic(clinic models.Clinic) (models.Clinic, error)
	DeleteClinic(id uint) error
}

type clinicService struct {
	clinicRepository clinicRepository.ClinicRepository
}

func NewClinicService(clinicRepository clinicRepository.ClinicRepository) *clinicService {
	return &clinicService{
		clinicRepository: clinicRepository,
	}
}

func (s *clinicService) GetClinics() ([]models.Clinic, error) {
	return s.clinicRepository.GetClinics()
}

func (s *clinicService) GetClinic(id uint) (models.Clinic, error) {
	return s.clinicRepository.GetClinic(id)
}

func (s *clinicService) CreateClinic(clinic models.Clinic) (models.Clinic, error) {
	return s.clinicRepository.CreateClinic(clinic)
}

func (s *clinicService) UpdateClinic(clinic models.Clinic) (models.Clinic, error) {
	return s.clinicRepository.UpdateClinic(clinic)
}

func (s *clinicService) DeleteClinic(id uint) error {
	return s.clinicRepository.DeleteClinic(id)
}
