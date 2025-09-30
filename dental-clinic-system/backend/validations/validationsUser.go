package validations

import (
	"dental-clinic-system/models/user"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func UserValidation(user *user.User) error {
	err := UserNamesValidation(user)
	if err != nil {
		return err
	}

	err = UserEmailValidation(user)
	if err != nil {
		return err
	}

	err = UserPasswordValidation(user)
	if err != nil {
		return err
	}

	err = ValidateUserPhones(user)
	if err != nil {
		return err
	}

	err = ValidateUserNationalID(user)
	if err != nil {
		return err
	}

	return nil
}

func UserNamesValidation(user *user.User) error {
	if user.FirstName == "" || user.LastName == "" {
		return errors.New("first name or last name is can not be empty")
	}

	if len(user.FirstName) < 2 || len(user.FirstName) > 50 || len(user.LastName) < 2 || len(user.LastName) > 50 {
		return errors.New("first name and last name must be between 2 and 50 characters")
	}

	nameRegex := `^[a-zA-ZçÇğĞıİöÖşŞüÜ]+$`
	if !regexp.MustCompile(nameRegex).MatchString(user.FirstName) || !regexp.MustCompile(nameRegex).MatchString(user.LastName) {
		return errors.New("first name and last name can only contain alphabetic characters")
	}

	caser := cases.Title(language.Und)
	user.FirstName = caser.String(strings.ToLower(user.FirstName))
	user.LastName = caser.String(strings.ToLower(user.LastName))

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	return nil

}

func UserEmailValidation(user *user.User) error {
	if user.Email == "" {
		return errors.New("email can not be empty")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(user.Email) {
		return errors.New("email is not valid")
	}

	user.Email = strings.TrimSpace(user.Email)

	return nil
}

func UserPasswordValidation(user *user.User) error {
	if user.Password == "" {
		return errors.New("password can not be empty")
	}

	if strings.Contains(user.Password, " ") {
		return errors.New("password can not contain spaces")
	}

	if len(user.Password) < 6 || len(user.Password) > 50 {
		return errors.New("password must be between 6 and 50 characters")
	}

	return nil
}

func ValidateUserPhones(user *user.User) error {
	countryCodePattern := `^\+?[0-9]{1,3}$`
	if !regexp.MustCompile(countryCodePattern).MatchString(user.CountryCode) {
		return errors.New("country code must be 1 to 3 digits, and can start with '+'")
	}

	personalPhonePattern := `^[0-9]{10}$`
	if !regexp.MustCompile(personalPhonePattern).MatchString(user.PhoneNumber) {
		return errors.New("personal phone number must be exactly 10 digits")
	}

	return nil
}

func ValidateUserNationalID(user *user.User) error {
	nationalIDPattern := `^[1-9]{1}[0-9]{9}[02468]{1}$`
	if !regexp.MustCompile(nationalIDPattern).MatchString(user.NationalID) {
		return errors.New("national ID must be exactly 11 digits")
	}

	return nil
}
