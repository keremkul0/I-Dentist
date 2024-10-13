package loginRepository

import (
	"dental-clinic-system/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRepository interface {
	Login(email string, password string) (models.Login, error)
}

func NewLoginRepository(db *gorm.DB) *loginRepository {
	return &loginRepository{DB: db}
}

type loginRepository struct {
	DB *gorm.DB
}

func (r *loginRepository) Login(email string, password string) (models.Login, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return models.Login{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.Login{}, err
	}

	return models.Login{
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
