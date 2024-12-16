package procedureService

import (
	"context"
	"dental-clinic-system/models"
)

type ProcedureRepository interface {
	GetProcedures(ctx context.Context, ClinicID uint) ([]models.Procedure, error)
	GetProcedure(ctx context.Context, id uint) (models.Procedure, error)
	CreateProcedure(ctx context.Context, procedure models.Procedure) (models.Procedure, error)
	UpdateProcedure(ctx context.Context, procedure models.Procedure) (models.Procedure, error)
	DeleteProcedure(ctx context.Context, id uint) error
}

type procedureService struct {
	procedureRepository ProcedureRepository
}

func NewProcedureService(procedureRepository ProcedureRepository) *procedureService {
	return &procedureService{
		procedureRepository: procedureRepository,
	}
}

func (s *procedureService) GetProcedures(ctx context.Context, ClinicID uint) ([]models.Procedure, error) {
	return s.procedureRepository.GetProcedures(ctx, ClinicID)
}

func (s *procedureService) GetProcedure(ctx context.Context, id uint) (models.Procedure, error) {
	return s.procedureRepository.GetProcedure(ctx, id)
}

func (s *procedureService) CreateProcedure(ctx context.Context, procedure models.Procedure) (models.Procedure, error) {
	return s.procedureRepository.CreateProcedure(ctx, procedure)
}

func (s *procedureService) UpdateProcedure(ctx context.Context, procedure models.Procedure) (models.Procedure, error) {
	return s.procedureRepository.UpdateProcedure(ctx, procedure)
}

func (s *procedureService) DeleteProcedure(ctx context.Context, id uint) error {
	return s.procedureRepository.DeleteProcedure(ctx, id)
}
