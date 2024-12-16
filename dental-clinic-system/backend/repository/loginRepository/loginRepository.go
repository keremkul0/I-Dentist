package loginRepository

import (
	"context"
	"dental-clinic-system/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewLoginRepository(db *gorm.DB) *loginRepository {
	return &loginRepository{DB: db}
}

type loginRepository struct {
	DB *gorm.DB
}

func (r *loginRepository) Login(ctx context.Context, email string, password string) (models.Login, error) {
	var user models.User
	result := r.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.Login{}, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return models.Login{
		Email:    user.Email,
		Password: user.Password,
	}, err
}
