package userGetModelRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type UserGetModelRepository interface {
	GetUserGetModels() ([]models.UserGetModel, error)
	GetUserGetModel(id uint) (models.UserGetModel, error)
	CreateUserGetModel(user models.UserGetModel) (models.UserGetModel, error)
	UpdateUserGetModel(user models.UserGetModel) (models.UserGetModel, error)
	DeleteUserGetModel(id uint) error
	CheckUserGetModelExist(user models.UserGetModel) bool
}

func NewUserGetModelRepository(db *gorm.DB) *userGetModelRepository {
	return &userGetModelRepository{DB: db}
}

type userGetModelRepository struct {
	DB *gorm.DB
}

func (r *userGetModelRepository) GetUserGetModels() ([]models.UserGetModel, error) {
	var userGetModels []models.UserGetModel
	if result := r.DB.Find(&userGetModels); result.Error != nil {
		return nil, result.Error
	}
	return userGetModels, nil
}

func (r *userGetModelRepository) GetUserGetModel(id uint) (models.UserGetModel, error) {
	var userGetModel models.UserGetModel
	if result := r.DB.First(&userGetModel, id); result.Error != nil {
		return models.UserGetModel{}, result.Error
	}
	return userGetModel, nil
}

func (r *userGetModelRepository) CreateUserGetModel(user models.UserGetModel) (models.UserGetModel, error) {
	if result := r.DB.Create(&user); result.Error != nil {
		return models.UserGetModel{}, result.Error
	}
	return user, nil
}

func (r *userGetModelRepository) UpdateUserGetModel(user models.UserGetModel) (models.UserGetModel, error) {
	if result := r.DB.Save(&user); result.Error != nil {
		return models.UserGetModel{}, result.Error
	}
	return user, nil
}

func (r *userGetModelRepository) DeleteUserGetModel(id uint) error {
	if result := r.DB.Delete(&models.UserGetModel{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userGetModelRepository) CheckUserGetModelExist(user models.UserGetModel) bool {
	var count int64
	r.DB.Model(&models.User{}).Where("id = ? or email = ? OR phone_number = ? OR national_id = ?", user.ID, user.Email, user.PhoneNumber, user.NationalID).First(&models.User{}).Count(&count)
	return count > 0
}
