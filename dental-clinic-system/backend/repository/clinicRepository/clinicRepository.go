package clinicRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type ClinicRepository interface {
	GetClinics() ([]models.Clinic, error)
	GetClinic(id uint) (models.Clinic, error)
	CreateClinic(clinic models.Clinic) (models.Clinic, error)
	UpdateClinic(clinic models.Clinic) (models.Clinic, error)
	DeleteClinic(id uint) error
	CheckClinicExist(clinic models.Clinic) bool
}

func NewClinicRepository(db *gorm.DB) *clinicRepository {
	return &clinicRepository{DB: db}
}

type clinicRepository struct {
	DB *gorm.DB
}

func (r *clinicRepository) GetClinics() ([]models.Clinic, error) {
	var clinics []models.Clinic
	err := r.DB.Find(&clinics).Error
	return clinics, err
}

func (r *clinicRepository) GetClinic(id uint) (models.Clinic, error) {
	var clinic models.Clinic
	err := r.DB.First(&clinic, id).Error
	return clinic, err
}

func (r *clinicRepository) CreateClinic(clinic models.Clinic) (models.Clinic, error) {

	err := r.DB.Create(&clinic).Error
	return clinic, err
}
func (r *clinicRepository) UpdateClinic(clinic models.Clinic) (models.Clinic, error) {
	err := r.DB.Save(&clinic).Error
	return clinic, err
}

func (r *clinicRepository) DeleteClinic(id uint) error {
	err := r.DB.Delete(&models.Clinic{}, id).Error
	return err
}

func (r *clinicRepository) CheckClinicExist(clinic models.Clinic) bool {
	var count int64
	r.DB.Model(&models.Clinic{}).Where("id = ? or email = ? or name = ? OR phone_number = ?", clinic.ID, clinic.Email, clinic.Name, clinic.PhoneNumber).First(&models.Clinic{}).Count(&count)
	return count > 0
}
