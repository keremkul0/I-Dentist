package authRepository

import (
	"dental-clinic-system/models"
	"errors"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(email string, password string) (models.Auth, error)
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{DB: db}
}

type authRepository struct {
	DB *gorm.DB
}

func (r *authRepository) Login(email string, password string) (models.Auth, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return models.Auth{}, err
	}

	//Compare the password from the request with the password from the database
	if user.Password != password {
		return models.Auth{}, errors.New("Invalid email or password")
	}

	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	//	return models.Auth{}, err
	//}

	return models.Auth{
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
