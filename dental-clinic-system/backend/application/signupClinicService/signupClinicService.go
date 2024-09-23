package signupClinicService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/clinicRepository"
	"dental-clinic-system/repository/userRepository"
)

type SignUpClinicService interface {
	SignUpClinic(clinic models.Clinic, user models.User) (models.Clinic, models.User, error)
}

type signUpClinicService struct {
	clinicRepository clinicRepository.ClinicRepository
	userRepository   userRepository.UserRepository
}

func NewSignUpClinicService(clinicRepository clinicRepository.ClinicRepository,
	userRepository userRepository.UserRepository) *signUpClinicService {
	return &signUpClinicService{
		clinicRepository: clinicRepository,
		userRepository:   userRepository,
	}
}

func (s *signUpClinicService) SignUpClinic(clinic models.Clinic, user models.User) (models.Clinic, models.User, error) {
	if s.userRepository.CheckUserExist(user) {
		return models.Clinic{}, models.User{}, nil
	}
	if s.clinicRepository.CheckClinicExist(clinic) {
		return models.Clinic{}, models.User{}, nil
	}

	clinic, err := s.clinicRepository.CreateClinic(clinic)

	if err != nil {
		return models.Clinic{}, models.User{}, err
	}
	user.ClinicID = clinic.ID
	user, err = s.userRepository.CreateUser(user)
	if err != nil {
		return models.Clinic{}, models.User{}, err
	}
	return clinic, user, nil
}
