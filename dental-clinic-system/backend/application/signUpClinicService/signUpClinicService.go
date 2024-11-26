package signUpClinicService

import (
	"dental-clinic-system/mapper"
	"dental-clinic-system/models"
	"dental-clinic-system/validations"
	"errors"
)

type UserRepository interface {
	CreateUserRepo(user models.User) (models.User, error)
	CheckUserExistRepo(user models.UserGetModel) bool
}

type ClinicRepository interface {
	CreateClinic(clinic models.Clinic) (models.Clinic, error)
	CheckClinicExist(clinic models.Clinic) bool
}

type CacheUserRepository interface {
	GetUserRepo(cacheKey string) (models.User, error)
}

type signUpClinicService struct {
	clinicRepository ClinicRepository
	userRepository   UserRepository
	repository       CacheUserRepository
}

func NewSignUpClinicService(clinicRepository ClinicRepository,
	userRepository UserRepository, repository CacheUserRepository) *signUpClinicService {
	return &signUpClinicService{
		clinicRepository: clinicRepository,
		userRepository:   userRepository,
		repository:       repository,
	}
}

func (s *signUpClinicService) SignUpClinic(clinic models.Clinic, userCacheKey string) (models.Clinic, models.UserGetModel, error) {

	user, err := s.repository.GetUserRepo(userCacheKey)

	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, errors.New("user not found")
	}

	userGetModel := mapper.UserMapper(user)
	if s.userRepository.CheckUserExistRepo(userGetModel) {
		return models.Clinic{}, models.UserGetModel{}, errors.New("user already exist")
	}

	if s.clinicRepository.CheckClinicExist(clinic) {
		return models.Clinic{}, models.UserGetModel{}, errors.New("clinic already exist")
	}

	err = validations.ClinicValidation(&clinic)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, errors.New("clinic validation error")
	}

	clinic, err = s.clinicRepository.CreateClinic(clinic)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, err
	}

	user.ClinicID = clinic.ID
	userGetModel.ClinicID = user.ClinicID

	user, err = s.userRepository.CreateUserRepo(user)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, err
	}

	return clinic, userGetModel, nil
}