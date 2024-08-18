package appointment

import (
	"dental-clinic-system/repository/appointment"
	"dental-clinic-system/repository/models"
)

type AppointmentService interface {
	GetAppointments() ([]models.Appointment, error)
	GetAppointment(id string) (models.Appointment, error)
	CreateAppointment(appointment models.Appointment) (models.Appointment, error)
	UpdateAppointment(appointment models.Appointment) (models.Appointment, error)
	DeleteAppointment(id string) error
}

type appointmentService struct {
	appointmentRepository appointment.AppointmentRepository
}

func NewAppointmentService(appointmentRepository appointment.AppointmentRepository) *appointmentService {
	return &appointmentService{
		appointmentRepository: appointmentRepository,
	}
}

func (s *appointmentService) GetAppointments() ([]models.Appointment, error) {
	return s.appointmentRepository.GetAppointments()
}

func (s *appointmentService) GetAppointment(id string) (models.Appointment, error) {
	return s.appointmentRepository.GetAppointment(id)
}

func (s *appointmentService) CreateAppointment(appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.CreateAppointment(appointment)
}

func (s *appointmentService) UpdateAppointment(appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.UpdateAppointment(appointment)
}

func (s *appointmentService) DeleteAppointment(id string) error {
	return s.appointmentRepository.DeleteAppointment(id)
}
