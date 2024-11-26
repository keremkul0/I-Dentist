package validations

import (
	"dental-clinic-system/models"
	"errors"
	"regexp"
	"strings"
)

func ClinicValidation(clinic *models.Clinic) error {
	err := ClinicNameValidation(clinic)
	if err != nil {
		return err
	}

	err = ClinicAddressValidation(clinic)
	if err != nil {
		return err
	}

	err = ClinicPhoneNumberValidation(clinic)
	if err != nil {
		return err
	}

	err = ClinicEmailValidation(clinic)
	if err != nil {
		return err
	}

	return nil
}

func ClinicNameValidation(clinic *models.Clinic) error {

	if clinic.Name == "" {
		return errors.New("Clinic name can not be empty")
	}

	if len(clinic.Name) < 2 || len(clinic.Name) > 50 {
		return errors.New("Clinic name must be between 2 and 50 characters")
	}

	clinic.Name = strings.TrimSpace(clinic.Name)
	return nil
}

func ClinicAddressValidation(clinic *models.Clinic) error {

	if clinic.Address == "" {
		return errors.New("Clinic address can not be empty")
	}

	if len(clinic.Address) < 2 || len(clinic.Address) > 50 {
		return errors.New("Clinic address must be between 2 and 50 characters")
	}

	clinic.Address = strings.TrimSpace(clinic.Address)
	return nil
}

func ClinicPhoneNumberValidation(clinic *models.Clinic) error {

	if clinic.PhoneNumber == "" {
		return errors.New("Clinic phone number can not be empty")
	}

	personalPhonePattern := `^[0-9]{10}$`
	if !regexp.MustCompile(personalPhonePattern).MatchString(clinic.PhoneNumber) {
		return errors.New("Personal phone number must be exactly 10 digits")
	}

	clinic.PhoneNumber = strings.TrimSpace(clinic.PhoneNumber)
	return nil
}

func ClinicEmailValidation(clinic *models.Clinic) error {

	if clinic.Email == "" {
		return errors.New("Clinic email can not be empty")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(clinic.Email) {
		return errors.New("Email is not valid")
	}

	clinic.Email = strings.TrimSpace(clinic.Email)
	return nil
}
