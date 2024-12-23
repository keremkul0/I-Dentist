package signUpClinicService

import (
	"context"
	"dental-clinic-system/mapper"
	"dental-clinic-system/models"
	"dental-clinic-system/validations"
	"errors"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	CheckUserExist(ctx context.Context, user models.UserGetModel) bool
}

type ClinicRepository interface {
	CreateClinic(ctx context.Context, clinic models.Clinic) (models.Clinic, error)
	CheckClinicExist(ctx context.Context, clinic models.Clinic) bool
}

type RedisRepository interface {
	GetData(ctx context.Context, cacheKey string, target any) error
	DeleteData(ctx context.Context, ID string) error
}

type signUpClinicService struct {
	clinicRepository ClinicRepository
	userRepository   UserRepository
	redisRepository  RedisRepository
}

func NewSignUpClinicService(clinicRepository ClinicRepository, userRepository UserRepository, redisRepository RedisRepository) *signUpClinicService {
	return &signUpClinicService{
		clinicRepository: clinicRepository,
		userRepository:   userRepository,
		redisRepository:  redisRepository,
	}
}

func (s *signUpClinicService) SignUpClinic(ctx context.Context, clinic models.Clinic, userCacheKey string) (models.Clinic, models.UserGetModel, error) {
	var user models.User
	err := s.redisRepository.GetData(ctx, userCacheKey, &user)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, errors.New("user not found")
	}

	// user, ok := data.(models.User) yerine aşağıdaki gibi bir kontrol yapılmalı

	//var user models.User
	//err = json.Unmarshal(data, &user)
	//if err != nil {
	//	return models.Clinic{}, models.UserGetModel{}, errors.New("failed to unmarshal user data")
	//}

	userGetModel := mapper.UserMapper(user)
	if s.userRepository.CheckUserExist(ctx, userGetModel) {
		return models.Clinic{}, models.UserGetModel{}, errors.New("user already exist")
	}

	if s.clinicRepository.CheckClinicExist(ctx, clinic) {
		return models.Clinic{}, models.UserGetModel{}, errors.New("clinic already exist")
	}

	err = validations.ClinicValidation(&clinic)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, errors.New("clinic validation error")
	}

	clinic, err = s.clinicRepository.CreateClinic(ctx, clinic)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, err
	}

	user.ClinicID = clinic.ID
	userGetModel.ClinicID = user.ClinicID

	user, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, err
	}

	err = s.redisRepository.DeleteData(ctx, userCacheKey)
	if err != nil {
		log.Warn().Err(err).Msg("failed to delete user cache")
	}

	return clinic, userGetModel, nil
}
