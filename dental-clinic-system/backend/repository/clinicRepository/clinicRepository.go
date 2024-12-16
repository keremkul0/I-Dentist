package clinicRepository

import (
	"context"
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

func NewClinicRepository(db *gorm.DB) *clinicRepository {
	return &clinicRepository{DB: db}
}

type clinicRepository struct {
	DB *gorm.DB
}

func (r *clinicRepository) GetClinics(ctx context.Context) ([]models.Clinic, error) {
	var clinics []models.Clinic
	err := r.DB.WithContext(ctx).Find(&clinics).Error
	return clinics, err
}

func (r *clinicRepository) GetClinic(ctx context.Context, id uint) (models.Clinic, error) {
	var clinic models.Clinic
	err := r.DB.WithContext(ctx).First(&clinic, id).Error
	return clinic, err
}

func (r *clinicRepository) CreateClinic(ctx context.Context, clinic models.Clinic) (models.Clinic, error) {
	err := r.DB.WithContext(ctx).Create(&clinic).Error
	return clinic, err
}

func (r *clinicRepository) UpdateClinic(ctx context.Context, clinic models.Clinic) (models.Clinic, error) {
	err := r.DB.WithContext(ctx).Save(&clinic).Error
	return clinic, err
}

func (r *clinicRepository) DeleteClinic(ctx context.Context, id uint) error {
	err := r.DB.WithContext(ctx).Delete(&models.Clinic{}, id).Error
	return err
}

func (r *clinicRepository) CheckClinicExist(ctx context.Context, clinic models.Clinic) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&models.Clinic{}).Where("id = ? or email = ? or name = ? OR phone_number = ?", clinic.ID, clinic.Email, clinic.Name, clinic.PhoneNumber).Count(&count)
	return count > 0
}
