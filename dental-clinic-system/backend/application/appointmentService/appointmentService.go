package appointmentService

import (
	"context"
	"dental-clinic-system/models"
)

type AppointmentRepository interface {
	GetAppointmentsRepo(ctx context.Context, ClinicID uint) ([]models.Appointment, error)
	GetAppointmentRepo(ctx context.Context, id uint) (models.Appointment, error)
	CreateAppointmentRepo(ctx context.Context, appointment models.Appointment) (models.Appointment, error)
	UpdateAppointmentRepo(ctx context.Context, appointment models.Appointment) (models.Appointment, error)
	DeleteAppointmentRepo(ctx context.Context, id uint) error
	GetDoctorAppointmentsRepo(ctx context.Context, id uint) ([]models.Appointment, error)
	GetPatientAppointmentsRepo(ctx context.Context, id uint) ([]models.Appointment, error)
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
	return s.appointmentRepository.GetAppointmentsRepo(ctx, ClinicID)
}

func (s *appointmentService) GetAppointment(ctx context.Context, id uint) (models.Appointment, error) {
	return s.appointmentRepository.GetAppointmentRepo(ctx, id)
}

func (s *appointmentService) CreateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.CreateAppointmentRepo(ctx, appointment)
}

func (s *appointmentService) UpdateAppointment(ctx context.Context, appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.UpdateAppointmentRepo(ctx, appointment)
}

func (s *appointmentService) DeleteAppointment(ctx context.Context, id uint) error {
	return s.appointmentRepository.DeleteAppointmentRepo(ctx, id)
}

func (s *appointmentService) GetDoctorAppointments(ctx context.Context, id uint) ([]models.Appointment, error) {
	return s.appointmentRepository.GetDoctorAppointmentsRepo(ctx, id)
}

func (s *appointmentService) GetPatientAppointments(ctx context.Context, id uint) ([]models.Appointment, error) {
	return s.appointmentRepository.GetPatientAppointmentsRepo(ctx, id)
}
