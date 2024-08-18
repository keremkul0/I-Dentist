package appointment

import (
	models2 "dental-clinic-system/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetAppointments() ([]models2.Appointment, error)
	GetAppointment(id string) (models2.Appointment, error)
	CreateAppointment(appointment models2.Appointment) (models2.Appointment, error)
	UpdateAppointment(appointment models2.Appointment) (models2.Appointment, error)
	GetAppointmentClinic(id string) (models2.Clinic, error)
	DeleteAppointment(id string) error
}

func NewAppointmentRepository(db *gorm.DB) *appointmentRepository {
	return &appointmentRepository{DB: db}
}

type appointmentRepository struct {
	DB *gorm.DB
}

func (r *appointmentRepository) GetAppointments() ([]models2.Appointment, error) {
	var appointments []models2.Appointment
	if result := r.DB.Preload("Clinic").Preload("Patient").Preload("Doctor").Find(&appointments); result.Error != nil {
		return nil, result.Error
	}
	return appointments, nil
}

func (r *appointmentRepository) GetAppointment(id string) (models2.Appointment, error) {
	var appointment models2.Appointment
	if result := r.DB.Preload("Clinic").Preload("Patient").Preload("Doctor").First(&appointment, id); result.Error != nil {
		return models2.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) CreateAppointment(appointment models2.Appointment) (models2.Appointment, error) {
	if result := r.DB.Create(&appointment); result.Error != nil {
		return models2.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) UpdateAppointment(appointment models2.Appointment) (models2.Appointment, error) {
	if result := r.DB.Save(&appointment); result.Error != nil {
		return models2.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) DeleteAppointment(id string) error {
	if result := r.DB.Delete(&models2.Appointment{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *appointmentRepository) GetAppointmentClinic(id string) (models2.Clinic, error) {

	var clinic models2.Clinic
	if result := r.DB.Preload("Appointments").First(&clinic, id); result.Error != nil {
		return models2.Clinic{}, result.Error
	}
	return clinic, nil
}
