package appointment

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetAppointments() ([]models.Appointment, error)
	GetAppointment(id string) (models.Appointment, error)
	CreateAppointment(appointment models.Appointment) (models.Appointment, error)
	UpdateAppointment(appointment models.Appointment) (models.Appointment, error)
	GetAppointmentClinic(id string) (models.Clinic, error)
	DeleteAppointment(id string) error
}

func NewAppointmentRepository(db *gorm.DB) *appointmentRepository {
	return &appointmentRepository{DB: db}
}

type appointmentRepository struct {
	DB *gorm.DB
}

func (r *appointmentRepository) GetAppointments() ([]models.Appointment, error) {
	var appointments []models.Appointment
	if result := r.DB.Preload("Clinic").Preload("Patient").Preload("Doctor").Find(&appointments); result.Error != nil {
		return nil, result.Error
	}
	return appointments, nil
}

func (r *appointmentRepository) GetAppointment(id string) (models.Appointment, error) {
	var appointment models.Appointment
	if result := r.DB.Preload("Clinic").Preload("Patient").Preload("Doctor").First(&appointment, id); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) CreateAppointment(appointment models.Appointment) (models.Appointment, error) {
	if result := r.DB.Create(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) UpdateAppointment(appointment models.Appointment) (models.Appointment, error) {
	if result := r.DB.Save(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) DeleteAppointment(id string) error {
	if result := r.DB.Delete(&models.Appointment{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *appointmentRepository) GetAppointmentClinic(id string) (models.Clinic, error) {

	var clinic models.Clinic
	if result := r.DB.Preload("Appointments").First(&clinic, id); result.Error != nil {
		return models.Clinic{}, result.Error
	}
	return clinic, nil
}
