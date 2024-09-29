package appointmentRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	GetAppointmentsRepo(ClinicID uint) ([]models.Appointment, error)
	GetAppointmentRepo(id uint) (models.Appointment, error)
	CreateAppointmentRepo(appointment models.Appointment) (models.Appointment, error)
	UpdateAppointmentRepo(appointment models.Appointment) (models.Appointment, error)
	DeleteAppointmentRepo(id uint) error
	GetDoctorAppointmentsRepo(id uint) ([]models.Appointment, error)
	GetPatientAppointmentsRepo(id uint) ([]models.Appointment, error)
}

func NewAppointmentRepository(db *gorm.DB) *appointmentRepository {
	return &appointmentRepository{DB: db}
}

type appointmentRepository struct {
	DB *gorm.DB
}

func (r *appointmentRepository) GetAppointmentsRepo(ClinicID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Where("clinic_id = ?", ClinicID).Preload("Clinic").Preload("Patient").Preload("Doctor").Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) GetAppointmentRepo(id uint) (models.Appointment, error) {
	var appointment models.Appointment
	if result := r.DB.Preload("Clinic").Preload("Patient").Preload("Doctor").First(&appointment, id); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) CreateAppointmentRepo(appointment models.Appointment) (models.Appointment, error) {
	if result := r.DB.Create(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) UpdateAppointmentRepo(appointment models.Appointment) (models.Appointment, error) {
	if result := r.DB.Save(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) DeleteAppointmentRepo(id uint) error {
	if result := r.DB.Delete(&models.Appointment{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *appointmentRepository) GetDoctorAppointmentsRepo(id uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Where("doctor_id = ?", id).Find(&appointments).Error
	return appointments, err

}

func (r *appointmentRepository) GetPatientAppointmentsRepo(id uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.Where("patient_id = ?", id).Find(&appointments).Error
	return appointments, err
}
