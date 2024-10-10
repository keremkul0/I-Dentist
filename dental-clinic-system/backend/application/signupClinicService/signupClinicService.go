package signupClinicService

import (
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"dental-clinic-system/repository/clinicRepository"
	"dental-clinic-system/repository/userRepository"
)

type SignUpClinicService interface {
	SignUpClinic(clinic models.Clinic, user models.User) (models.Clinic, models.UserGetModel, error)
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

func (s *signUpClinicService) SignUpClinic(clinic models.Clinic, user models.User) (models.Clinic, models.UserGetModel, error) {

	userGetModel := helpers.UserConvertor(user)

	if s.userRepository.CheckUserExistRepo(userGetModel) {
		return models.Clinic{}, models.UserGetModel{}, nil
	}
	if s.clinicRepository.CheckClinicExist(clinic) {
		return models.Clinic{}, models.UserGetModel{}, nil
	}

	clinic, err := s.clinicRepository.CreateClinic(clinic)

	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, err
	}

	user.ClinicID = clinic.ID
	user, err = s.userRepository.CreateUserRepo(user)
	userGetModel = helpers.UserConvertor(user)

	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, err
	}
	return clinic, userGetModel, nil
}
