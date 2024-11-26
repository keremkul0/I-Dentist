package appointmentService

import (
	"dental-clinic-system/models"
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

type appointmentService struct {
	appointmentRepository AppointmentRepository
}

func NewAppointmentService(appointmentRepository AppointmentRepository) *appointmentService {
	return &appointmentService{
		appointmentRepository: appointmentRepository,
	}
}

func (s *appointmentService) GetAppointments(ClinicID uint) ([]models.Appointment, error) {
	return s.appointmentRepository.GetAppointmentsRepo(ClinicID)
}

func (s *appointmentService) GetAppointment(id uint) (models.Appointment, error) {
	return s.appointmentRepository.GetAppointmentRepo(id)
}

func (s *appointmentService) CreateAppointment(appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.CreateAppointmentRepo(appointment)
}

func (s *appointmentService) UpdateAppointment(appointment models.Appointment) (models.Appointment, error) {
	return s.appointmentRepository.UpdateAppointmentRepo(appointment)
}

func (s *appointmentService) DeleteAppointment(id uint) error {
	return s.appointmentRepository.DeleteAppointmentRepo(id)
}

func (s *appointmentService) GetDoctorAppointments(id uint) ([]models.Appointment, error) {
	return s.appointmentRepository.GetDoctorAppointmentsRepo(id)
}

func (s *appointmentService) GetPatientAppointments(id uint) ([]models.Appointment, error) {
	return s.appointmentRepository.GetPatientAppointmentsRepo(id)
}
