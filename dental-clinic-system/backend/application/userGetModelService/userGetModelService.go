package userGetModelService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/userGetModelRepository"
)

type UserGetModelService interface {
	GetUserGetModels() ([]models.UserGetModel, error)
	GetUserGetModel(id uint) (models.UserGetModel, error)
	CreateUserGetModel(user models.UserGetModel) (models.UserGetModel, error)
	UpdateUserGetModel(user models.UserGetModel) (models.UserGetModel, error)
	DeleteUserGetModel(id uint) error
}

type userGetModelService struct {
	userGetModelRepository userGetModelRepository.UserGetModelRepository
}

func NewUserGetModelService(userGetModelRepository userGetModelRepository.UserGetModelRepository) *userGetModelService {
	return &userGetModelService{
		userGetModelRepository: userGetModelRepository,
	}
}

func (s *userGetModelService) GetUserGetModels() ([]models.UserGetModel, error) {
	return s.userGetModelRepository.GetUserGetModels()
}

func (s *userGetModelService) GetUserGetModel(id uint) (models.UserGetModel, error) {
	return s.userGetModelRepository.GetUserGetModel(id)
}

func (s *userGetModelService) CreateUserGetModel(user models.UserGetModel) (models.UserGetModel, error) {
	return s.userGetModelRepository.CreateUserGetModel(user)
}

func (s *userGetModelService) UpdateUserGetModel(user models.UserGetModel) (models.UserGetModel, error) {
	return s.userGetModelRepository.UpdateUserGetModel(user)
}

func (s *userGetModelService) DeleteUserGetModel(id uint) error {
	return s.userGetModelRepository.DeleteUserGetModel(id)
}
