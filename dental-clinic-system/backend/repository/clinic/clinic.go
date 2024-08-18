package clinic

import (
	models2 "dental-clinic-system/models"
	"gorm.io/gorm"
)

type ClinicRepository interface {
	GetClinics() ([]models2.Clinic, error)
	GetClinic(id uint) (models2.Clinic, error)
	CreateClinic(clinic models2.Clinic) (models2.Clinic, error)
	UpdateClinic(clinic models2.Clinic) (models2.Clinic, error)
	GetClinicAppointments(id uint) ([]models2.Appointment, error)
	DeleteClinic(id uint) error
}

func NewClinicRepository(db *gorm.DB) *clinicRepository {
	return &clinicRepository{DB: db}
}

type clinicRepository struct {
	DB *gorm.DB
}

func (r *clinicRepository) GetClinics() ([]models2.Clinic, error) {
	var clinics []models2.Clinic
	err := r.DB.Find(&clinics).Error
	return clinics, err
}

func (r *clinicRepository) GetClinic(id uint) (models2.Clinic, error) {
	var clinic models2.Clinic
	err := r.DB.First(&clinic, id).Error
	return clinic, err
}

func (r *clinicRepository) CreateClinic(clinic models2.Clinic) (models2.Clinic, error) {
	err := r.DB.Create(&clinic).Error
	return clinic, err
}

func (r *clinicRepository) UpdateClinic(clinic models2.Clinic) (models2.Clinic, error) {
	err := r.DB.Save(&clinic).Error
	return clinic, err
}

func (r *clinicRepository) GetClinicAppointments(id uint) ([]models2.Appointment, error) {
	var appointments []models2.Appointment
	err := r.DB.Where("clinic_id = ?", id).Find(&appointments).Error
	return appointments, err
}

func (r *clinicRepository) DeleteClinic(id uint) error {
	err := r.DB.Delete(&models2.Clinic{}, id).Error
	return err
}
