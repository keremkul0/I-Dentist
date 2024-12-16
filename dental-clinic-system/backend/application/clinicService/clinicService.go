package clinicService

import (
	"context"
	"dental-clinic-system/models"
)

type ClinicRepository interface {
	GetClinics(ctx context.Context) ([]models.Clinic, error)
	GetClinic(ctx context.Context, id uint) (models.Clinic, error)
	CreateClinic(ctx context.Context, clinic models.Clinic) (models.Clinic, error)
	UpdateClinic(ctx context.Context, clinic models.Clinic) (models.Clinic, error)
	DeleteClinic(ctx context.Context, id uint) error
}

type clinicService struct {
	clinicRepository ClinicRepository
}

func NewClinicService(clinicRepository ClinicRepository) *clinicService {
	return &clinicService{
		clinicRepository: clinicRepository,
	}
}

func (s *clinicService) GetClinics(ctx context.Context) ([]models.Clinic, error) {
	return s.clinicRepository.GetClinics(ctx)
}

func (s *clinicService) GetClinic(ctx context.Context, id uint) (models.Clinic, error) {
	return s.clinicRepository.GetClinic(ctx, id)
}

func (s *clinicService) CreateClinic(ctx context.Context, clinic models.Clinic) (models.Clinic, error) {
	return s.clinicRepository.CreateClinic(ctx, clinic)
}

func (s *clinicService) UpdateClinic(ctx context.Context, clinic models.Clinic) (models.Clinic, error) {
	return s.clinicRepository.UpdateClinic(ctx, clinic)
}

func (s *clinicService) DeleteClinic(ctx context.Context, id uint) error {
	return s.clinicRepository.DeleteClinic(ctx, id)
}
