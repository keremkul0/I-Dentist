package userRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUsersRepo(ClinicID uint) ([]models.User, error)
	GetUserRepo(id uint) (models.User, error)
	GetUserByEmailRepo(email string) (models.User, error)
	CreateUserRepo(user models.User) (models.User, error)
	UpdateUserRepo(user models.User) (models.User, error)
	DeleteUserRepo(id uint) error
	CheckUserExistRepo(user models.UserGetModel) bool
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{DB: db}
}

type userRepository struct {
	DB *gorm.DB
}

func (r *userRepository) GetUsersRepo(ClinicID uint) ([]models.User, error) {
	var users []models.User

	result := r.DB.Where("clinic_id = ?", ClinicID).Find(&users)
	if result.Error != nil {
		return []models.User{}, result.Error
	}

	return users, nil
}

func (r *userRepository) GetUserRepo(id uint) (models.User, error) {
	var user models.User

	result := r.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func (r *userRepository) GetUserByEmailRepo(email string) (models.User, error) {
	var user models.User

	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func (r *userRepository) CreateUserRepo(user models.User) (models.User, error) {
	r.DB.Create(&user)
	return user, nil
}

func (r *userRepository) UpdateUserRepo(user models.User) (models.User, error) {
	r.DB.Save(&user)
	return user, nil
}

func (r *userRepository) DeleteUserRepo(id uint) error {
	r.DB.Delete(&models.User{}, id)
	return nil
}

func (r *userRepository) CheckUserExistRepo(user models.UserGetModel) bool {
	var count int64
	r.DB.Model(&models.User{}).Where("national_id = ? OR email = ? OR phone_number = ?", user.NationalID, user.Email, user.PhoneNumber).Count(&count)
	return count > 0
}
