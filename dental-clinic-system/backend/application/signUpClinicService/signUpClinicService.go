package signUpClinicService

import (
	"context"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"dental-clinic-system/redisService"
	"dental-clinic-system/repository/clinicRepository"
	"dental-clinic-system/repository/userRepository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type SignUpClinicService interface {
	SignUpClinic(clinic models.Clinic, userCacheKey string) (models.Clinic, models.UserGetModel, error)
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

func (s *signUpClinicService) SignUpClinic(clinic models.Clinic, userCacheKey string) (models.Clinic, models.UserGetModel, error) {

	ctx := context.Background()
	val, err := redisService.Rdb.Get(ctx, userCacheKey).Result()
	if errors.Is(err, redis.Nil) {
		fmt.Println("Key does not exist")
	} else if err != nil {
		fmt.Println("Error retrieving data:", err)
	}

	var user models.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return models.Clinic{}, models.UserGetModel{}, fmt.Errorf("error unmarshalling user data: %v", err)
	}

	userGetModel := helpers.UserConvertor(user)
	if s.userRepository.CheckUserExistRepo(userGetModel) {
		return models.Clinic{}, models.UserGetModel{}, errors.New("user already exist")
	}

	if s.clinicRepository.CheckClinicExist(clinic) {
		return models.Clinic{}, models.UserGetModel{}, errors.New("clinic already exist")
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
