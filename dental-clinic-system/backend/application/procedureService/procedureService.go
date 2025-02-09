package procedureService

import (
	"context"
	"dental-clinic-system/models/procedure"
)

type ProcedureRepository interface {
	GetProcedures(ctx context.Context, ClinicID uint) ([]procedure.Procedure, error)
	GetProcedure(ctx context.Context, id uint) (procedure.Procedure, error)
	CreateProcedure(ctx context.Context, procedure procedure.Procedure) (procedure.Procedure, error)
	UpdateProcedure(ctx context.Context, procedure procedure.Procedure) (procedure.Procedure, error)
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

func (s *procedureService) GetProcedures(ctx context.Context, ClinicID uint) ([]procedure.Procedure, error) {
	return s.procedureRepository.GetProcedures(ctx, ClinicID)
}

func (s *procedureService) GetProcedure(ctx context.Context, id uint) (procedure.Procedure, error) {
	return s.procedureRepository.GetProcedure(ctx, id)
}

func (s *procedureService) CreateProcedure(ctx context.Context, procedure procedure.Procedure) (procedure.Procedure, error) {
	return s.procedureRepository.CreateProcedure(ctx, procedure)
}

func (s *procedureService) UpdateProcedure(ctx context.Context, procedure procedure.Procedure) (procedure.Procedure, error) {
	return s.procedureRepository.UpdateProcedure(ctx, procedure)
}

func (s *procedureService) DeleteProcedure(ctx context.Context, id uint) error {
	return s.procedureRepository.DeleteProcedure(ctx, id)
}
