package userRepository

import (
	"context"
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{DB: db}
}

type userRepository struct {
	DB *gorm.DB
}

func (r *userRepository) GetUsers(ctx context.Context, ClinicID uint) ([]models.User, error) {
	var users []models.User
	err := r.DB.WithContext(ctx).Where("clinic_id = ?", ClinicID).Find(&users).Error
	return users, err
}

func (r *userRepository) GetUser(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	err := r.DB.WithContext(ctx).Create(&user).Error
	return user, err
}

func (r *userRepository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	err := r.DB.WithContext(ctx).Save(&user).Error
	return user, err
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	err := r.DB.WithContext(ctx).Delete(&models.User{}, id).Error
	return err
}

func (r *userRepository) CheckUserExist(ctx context.Context, user models.UserGetModel) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&models.User{}).Where("national_id = ? OR email = ? OR phone_number = ?", user.NationalID, user.Email, user.PhoneNumber).Count(&count)
	return count > 0
}
