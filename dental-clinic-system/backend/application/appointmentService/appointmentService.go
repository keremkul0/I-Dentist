package appointmentService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/appointmentRepository"
)

type AppointmentService interface {
	GetAppointments() ([]models.Appointment, error)
	GetAppointment(id uint) (models.Appointment, error)
	CreateAppointment(appointment models.Appointment) (models.Appointment, error)
	UpdateAppointment(appointment models.Appointment) (models.Appointment, error)
	DeleteAppointment(id uint) error
}

type appointmentService struct {
	appointmentRepository appointmentRepository.AppointmentRepository
}

func NewAppointmentService(appointmentRepository appointmentRepository.AppointmentRepository) *appointmentService {
	return &appointmentService{
		appointmentRepository: appointmentRepository,
	}
}

func (s *appointmentService) GetAppointments() ([]models.Appointment, error) {
	return s.appointmentRepository.GetAppointments()
}

func (s *appointmentService) GetAppointment(id uint) (models.Appointment, error) {
	return s.appointmentRepository.GetAppointment(id)
}

func (s *appointmentService) CreateAppointment(appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.CreateAppointment(appointment)
}

func (s *appointmentService) UpdateAppointment(appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.UpdateAppointment(appointment)
}

func (s *appointmentService) DeleteAppointment(id uint) error {
	return s.appointmentRepository.DeleteAppointment(id)
}
