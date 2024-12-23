package procedureRepository

import (
	"context"
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

func NewProcedureRepository(db *gorm.DB) *procedureRepository {
	return &procedureRepository{DB: db}
}

type procedureRepository struct {
	DB *gorm.DB
}

func (r *procedureRepository) GetProcedures(ctx context.Context, ClinicID uint) ([]models.Procedure, error) {
	var procedures []models.Procedure
	result := r.DB.WithContext(ctx).Where("clinic_id = ?", ClinicID).Find(&procedures)
	return procedures, result.Error
}

func (r *procedureRepository) GetProcedure(ctx context.Context, id uint) (models.Procedure, error) {
	var procedure models.Procedure
	result := r.DB.WithContext(ctx).First(&procedure, id)
	return procedure, result.Error
}

func (r *procedureRepository) CreateProcedure(ctx context.Context, procedure models.Procedure) (models.Procedure, error) {
	result := r.DB.WithContext(ctx).Create(&procedure)
	return procedure, result.Error
}

func (r *procedureRepository) UpdateProcedure(ctx context.Context, procedure models.Procedure) (models.Procedure, error) {
	result := r.DB.WithContext(ctx).Save(&procedure)
	return procedure, result.Error
}

func (r *procedureRepository) DeleteProcedure(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.Procedure{}, id)
	return result.Error
}
