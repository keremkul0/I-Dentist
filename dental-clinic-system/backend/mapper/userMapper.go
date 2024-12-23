package mapper

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

// UserMapper is a struct that contains the methods to convert a user model to a user get model

func UserMapper(user models.User) models.UserGetModel {
	return models.UserGetModel{
		Model: gorm.Model{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: user.DeletedAt,
		},
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Roles:       user.Roles,
		ClinicID:    user.ClinicID,
		NationalID:  user.NationalID,
		LastLogin:   user.LastLogin,
		IsActive:    user.IsActive,
		PhoneNumber: user.PhoneNumber,
	}
}
