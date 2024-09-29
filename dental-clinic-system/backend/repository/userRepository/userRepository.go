package userRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUsersRepo(ClinicID uint) ([]models.UserGetModel, error)
	GetUserRepo(id uint) (models.UserGetModel, error)
	GetUserByEmailRepo(email string) (models.UserGetModel, error)
	CreateUserRepo(user models.User) (models.UserGetModel, error)
	UpdateUserRepo(user models.User) (models.UserGetModel, error)
	DeleteUserRepo(id uint) error
	CheckUserExistRepo(user models.UserGetModel) bool
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{DB: db}
}

type userRepository struct {
	DB *gorm.DB
}

func (r *userRepository) GetUsersRepo(ClinicID uint) ([]models.UserGetModel, error) {
	var users []models.UserGetModel
	r.DB.Where("clinic_id = ?", ClinicID).Find(&users)
	return users, nil
}

func (r *userRepository) GetUserRepo(id uint) (models.UserGetModel, error) {
	var user models.UserGetModel
	r.DB.First(&user, id)
	return user, nil
}

func (r *userRepository) GetUserByEmailRepo(email string) (models.UserGetModel, error) {
	var user models.UserGetModel
	r.DB.Where("email = ?", email).First(&user)
	return user, nil
}

func (r *userRepository) CreateUserRepo(user models.User) (models.UserGetModel, error) {
	r.DB.Create(&user)
	return r.GetUserRepo(user.ID)
}

func (r *userRepository) UpdateUserRepo(user models.User) (models.UserGetModel, error) {
	r.DB.Save(&user)
	return r.GetUserRepo(user.ID)
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
