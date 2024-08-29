package userRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUsers() ([]models.User, error)
	GetUser(id uint) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(id uint) error
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{DB: db}
}

type userRepository struct {
	DB *gorm.DB
}

func (r *userRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	if result := r.DB.Find(&users); result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) GetUser(id uint) (models.User, error) {
	var user models.User
	if result := r.DB.First(&user, id); result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	if result := r.DB.Create(&user); result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user models.User) (models.User, error) {
	if result := r.DB.Save(&user); result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) DeleteUser(id uint) error {
	if result := r.DB.Delete(&models.User{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}
