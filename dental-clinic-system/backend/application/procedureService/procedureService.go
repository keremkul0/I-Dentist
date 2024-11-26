package procedureService

import (
	"dental-clinic-system/models"
)

type ProcedureRepository interface {
	GetProcedures(ClinicID uint) ([]models.Procedure, error)
	GetProcedure(id uint) (models.Procedure, error)
	CreateProcedure(procedure models.Procedure) (models.Procedure, error)
	UpdateProcedure(procedure models.Procedure) (models.Procedure, error)
	DeleteProcedure(id uint) error
}

func NewProcedureService(procedureRepository ProcedureRepository) *procedureService {
	return &procedureService{
		procedureRepository: procedureRepository,
	}
}

type procedureService struct {
	procedureRepository ProcedureRepository
}

func (s *procedureService) GetProcedures(ClinicID uint) ([]models.Procedure, error) {
	return s.procedureRepository.GetProcedures(ClinicID)
}

func (s *procedureService) GetProcedure(id uint) (models.Procedure, error) {
	return s.procedureRepository.GetProcedure(id)
}

func (s *procedureService) CreateProcedure(procedure models.Procedure) (models.Procedure, error) {
	return s.procedureRepository.CreateProcedure(procedure)
}

func (s *procedureService) UpdateProcedure(procedure models.Procedure) (models.Procedure, error) {
	return s.procedureRepository.UpdateProcedure(procedure)
}

func (s *procedureService) DeleteProcedure(id uint) error {
	return s.procedureRepository.DeleteProcedure(id)
}
