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

type RedisRepository interface {
	GetData(key string) (any, error)
}

type signUpClinicService struct {
	clinicRepository ClinicRepository
	userRepository   UserRepository
	redisRepository  RedisRepository
}

func NewSignUpClinicService(clinicRepository ClinicRepository,
	userRepository UserRepository, redisRepository RedisRepository) *signUpClinicService {
	return &signUpClinicService{
		clinicRepository: clinicRepository,
		userRepository:   userRepository,
		redisRepository:  redisRepository,
	}
}

func (s *signUpClinicService) SignUpClinic(clinic models.Clinic, userCacheKey string) (models.Clinic, models.UserGetModel, error) {

	data, err := s.redisRepository.GetData(userCacheKey)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, errors.New("user not found")
	}

	user, ok := data.(models.User)
	if !ok {
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
