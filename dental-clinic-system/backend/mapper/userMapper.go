package mapper

import (
	"dental-clinic-system/models/user"

	"gorm.io/gorm"
)

// MapUserToUserGetModel is a struct that contains the methods to convert a user model to a user get model

// MapUserToUserGetModel converts a User entity to a UserGetModel for API responses.
func MapUserToUserGetModel(u user.User) user.UserGetModel {
	return user.UserGetModel{
		Model: gorm.Model{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: u.DeletedAt,
		},
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		Roles:       u.Roles,
		ClinicID:    u.ClinicID,
		NationalID:  u.NationalID,
		LastLogin:   u.LastLogin,
		IsActive:    u.IsActive,
		PhoneNumber: u.PhoneNumber,
	}
}
