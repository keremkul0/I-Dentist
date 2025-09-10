package clinic

import (
	"errors"

	"gorm.io/gorm"
)

type Clinic struct {
	gorm.Model
	Name        string `json:"name" gorm:"Index"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number" gorm:"uniqueIndex"`
	Email       string `json:"email" gorm:"uniqueIndex"`
}

// Error types
var (
	ErrClinicNotFound       = errors.New("clinic not found")
	ErrClinicAlreadyExists  = errors.New("clinic already exists")
	ErrClinicValidation     = errors.New("clinic validation errors")
	ErrClinicCreation       = errors.New("failed to create clinic")
	ErrClinicUpdate         = errors.New("failed to update clinic")
	ErrClinicDeletion       = errors.New("failed to delete clinic")
	ErrClinicExistenceCheck = errors.New("failed to check clinic existence")
)
