package appointmentService

import (
	"context"
	"dental-clinic-system/models"
)

type AppointmentRepository interface {
	GetAppointments(ctx context.Context, ClinicID uint) ([]models.Appointment, error)
	GetAppointment(ctx context.Context, id uint) (models.Appointment, error)
	CreateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error)
	UpdateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error)
	DeleteAppointment(ctx context.Context, id uint) error
	GetDoctorAppointments(ctx context.Context, id uint) ([]models.Appointment, error)
	GetPatientAppointments(ctx context.Context, id uint) ([]models.Appointment, error)
}

type appointmentService struct {
	appointmentRepository AppointmentRepository
}

func NewAppointmentService(appointmentRepository AppointmentRepository) *appointmentService {
	return &appointmentService{
		appointmentRepository: appointmentRepository,
	}
}

func (s *appointmentService) GetAppointments(ctx context.Context, ClinicID uint) ([]models.Appointment, error) {
	return s.appointmentRepository.GetAppointments(ctx, ClinicID)
}

func (s *appointmentService) GetAppointment(ctx context.Context, id uint) (models.Appointment, error) {
	return s.appointmentRepository.GetAppointment(ctx, id)
}

func (s *appointmentService) CreateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.CreateAppointment(ctx, appointment)
}

func (s *appointmentService) UpdateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.UpdateAppointment(ctx, appointment)
}

func (s *appointmentService) DeleteAppointment(ctx context.Context, id uint) error {
	return s.appointmentRepository.DeleteAppointment(ctx, id)
}

func (s *appointmentService) GetDoctorAppointments(ctx context.Context, id uint) ([]models.Appointment, error) {
	return s.appointmentRepository.GetDoctorAppointments(ctx, id)
}

func (s *appointmentService) GetPatientAppointments(ctx context.Context, id uint) ([]models.Appointment, error) {
	return s.appointmentRepository.GetPatientAppointments(ctx, id)
}
