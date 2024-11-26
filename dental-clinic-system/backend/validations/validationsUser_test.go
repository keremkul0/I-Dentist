package validations

import (
	"dental-clinic-system/models"
	"testing"
)

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name    string
		user    *models.User
		wantErr bool
	}{
		{
			name: "Valid user",
			user: &models.User{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Password:    "securepassword",
				CountryCode: "+1",
				PhoneNumber: "1234567890",
				NationalID:  "12345678901",
			},
			wantErr: false,
		},
		{
			name: "Empty first name",
			user: &models.User{
				FirstName:   "",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Password:    "securepassword",
				CountryCode: "+1",
				PhoneNumber: "1234567890",
				NationalID:  "12345678901",
			},
			wantErr: true,
		},
		{
			name: "Invalid email",
			user: &models.User{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "invalid-email",
				Password:    "securepassword",
				CountryCode: "+1",
				PhoneNumber: "1234567890",
				NationalID:  "12345678901",
			},
			wantErr: true,
		},
		{
			name: "Short password",
			user: &models.User{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Password:    "short",
				CountryCode: "+1",
				PhoneNumber: "1234567890",
				NationalID:  "12345678901",
			},
			wantErr: true,
		},
		{
			name: "Invalid phone number",
			user: &models.User{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Password:    "securepassword",
				CountryCode: "+1",
				PhoneNumber: "12345",
				NationalID:  "12345678901",
			},
			wantErr: true,
		},
		{
			name: "Invalid national ID",
			user: &models.User{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john.doe@example.com",
				Password:    "securepassword",
				CountryCode: "+1",
				PhoneNumber: "1234567890",
				NationalID:  "12345",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserValidation(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
