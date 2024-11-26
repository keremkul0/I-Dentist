package procedureRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

func NewProcedureRepository(db *gorm.DB) *procedureRepository {
	return &procedureRepository{DB: db}
}

type procedureRepository struct {
	DB *gorm.DB
}

func (r *procedureRepository) GetProcedures(ClinicID uint) ([]models.Procedure, error) {
	var procedures []models.Procedure
	if result := r.DB.Where("clinic_id = ?", ClinicID).Find(&procedures); result.Error != nil {
		return []models.Procedure{}, result.Error
	}
	return procedures, nil
}

func (r *procedureRepository) GetProcedure(id uint) (models.Procedure, error) {
	var procedure models.Procedure
	if result := r.DB.First(&procedure, id); result.Error != nil {
		return models.Procedure{}, result.Error
	}
	return procedure, nil
}

func (r *procedureRepository) CreateProcedure(procedure models.Procedure) (models.Procedure, error) {
	if result := r.DB.Create(&procedure); result.Error != nil {
		return models.Procedure{}, result.Error
	}
	return procedure, nil
}

func (r *procedureRepository) UpdateProcedure(procedure models.Procedure) (models.Procedure, error) {
	if result := r.DB.Save(&procedure); result.Error != nil {
		return models.Procedure{}, result.Error
	}
	return procedure, nil
}

func (r *procedureRepository) DeleteProcedure(id uint) error {
	if result := r.DB.Delete(&models.Procedure{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}
