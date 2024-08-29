package procedureService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/procedureRepository"
)

// ProcedureService describes the service.

type ProcedureService interface {
	GetProcedures() ([]models.Procedure, error)
	GetProcedure(id uint) (models.Procedure, error)
	CreateProcedure(procedure models.Procedure) (models.Procedure, error)
	UpdateProcedure(procedure models.Procedure) (models.Procedure, error)
	DeleteProcedure(id uint) error
}

// NewProcedureService creates a new procedure service.

func NewProcedureService(procedureRepository procedureRepository.ProcedureRepository) *procedureService {
	return &procedureService{
		procedureRepository: procedureRepository,
	}
}

type procedureService struct {
	procedureRepository procedureRepository.ProcedureRepository
}

// GetProcedures returns all procedures.

func (s *procedureService) GetProcedures() ([]models.Procedure, error) {
	return s.procedureRepository.GetProcedures()
}

// GetProcedure returns a procedure by its ID.

func (s *procedureService) GetProcedure(id uint) (models.Procedure, error) {
	return s.procedureRepository.GetProcedure(id)
}

// CreateProcedure creates a new procedure.

func (s *procedureService) CreateProcedure(procedure models.Procedure) (models.Procedure, error) {
	return s.procedureRepository.CreateProcedure(procedure)
}

// UpdateProcedure updates a procedure.

func (s *procedureService) UpdateProcedure(procedure models.Procedure) (models.Procedure, error) {
	return s.procedureRepository.UpdateProcedure(procedure)
}

// DeleteProcedure deletes a procedure by its ID.

func (s *procedureService) DeleteProcedure(id uint) error {
	return s.procedureRepository.DeleteProcedure(id)
}
