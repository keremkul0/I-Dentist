package validations

import (
	"dental-clinic-system/models/clinic"
	"errors"
	"regexp"
	"strings"
)

func ClinicValidation(clinic *clinic.Clinic) error {
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

func ClinicNameValidation(clinic *clinic.Clinic) error {

	if clinic.Name == "" {
		return errors.New("clinic name can not be empty")
	}

	if len(clinic.Name) < 2 || len(clinic.Name) > 50 {
		return errors.New("clinic name must be between 2 and 50 characters")
	}

	clinic.Name = strings.TrimSpace(clinic.Name)
	return nil
}

func ClinicAddressValidation(clinic *clinic.Clinic) error {

	if clinic.Address == "" {
		return errors.New("clinic address can not be empty")
	}

	if len(clinic.Address) < 2 || len(clinic.Address) > 50 {
		return errors.New("clinic address must be between 2 and 50 characters")
	}

	clinic.Address = strings.TrimSpace(clinic.Address)
	return nil
}

func ClinicPhoneNumberValidation(clinic *clinic.Clinic) error {

	if clinic.PhoneNumber == "" {
		return errors.New("clinic phone number can not be empty")
	}

	personalPhonePattern := `^[0-9]{10}$`
	if !regexp.MustCompile(personalPhonePattern).MatchString(clinic.PhoneNumber) {
		return errors.New("personal phone number must be exactly 10 digits")
	}

	clinic.PhoneNumber = strings.TrimSpace(clinic.PhoneNumber)
	return nil
}

func ClinicEmailValidation(clinic *clinic.Clinic) error {

	if clinic.Email == "" {
		return errors.New("clinic email can not be empty")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(clinic.Email) {
		return errors.New("email is not valid")
	}

	clinic.Email = strings.TrimSpace(clinic.Email)
	return nil
}
