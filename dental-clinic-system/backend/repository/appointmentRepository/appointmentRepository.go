package appointmentRepository

import (
	"context"
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

func NewAppointmentRepository(db *gorm.DB) *appointmentRepository {
	return &appointmentRepository{DB: db}
}

type appointmentRepository struct {
	DB *gorm.DB
}

func (r *appointmentRepository) GetAppointments(ctx context.Context, ClinicID uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.WithContext(ctx).Where("clinic_id = ?", ClinicID).Preload("Clinic").Preload("Patient").Preload("Doctor").Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) GetAppointment(ctx context.Context, id uint) (models.Appointment, error) {
	var appointment models.Appointment
	if result := r.DB.WithContext(ctx).Preload("Clinic").Preload("Patient").Preload("Doctor").First(&appointment, id); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) CreateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error) {
	if result := r.DB.WithContext(ctx).Create(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) UpdateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error) {
	if result := r.DB.WithContext(ctx).Save(&appointment); result.Error != nil {
		return models.Appointment{}, result.Error
	}
	return appointment, nil
}

func (r *appointmentRepository) DeleteAppointment(ctx context.Context, id uint) error {
	if result := r.DB.WithContext(ctx).Delete(&models.Appointment{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *appointmentRepository) GetDoctorAppointments(ctx context.Context, id uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.WithContext(ctx).Where("doctor_id = ?", id).Find(&appointments).Error
	return appointments, err

}

func (r *appointmentRepository) GetPatientAppointments(ctx context.Context, id uint) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.DB.WithContext(ctx).Where("patient_id = ?", id).Find(&appointments).Error
	return appointments, err
}
