package validations

import (
	"dental-clinic-system/models/clinic"
	"testing"
)

func TestClinicValidation(t *testing.T) {
	tests := []struct {
		name    string
		clinic  *clinic.Clinic
		wantErr bool
	}{
		{
			name: "Valid clinic",
			clinic: &clinic.Clinic{
				Name:        "Healthy Smiles",
				Address:     "123 Dental St",
				PhoneNumber: "1234567890",
				Email:       "contact@healthysmiles.com",
			},
			wantErr: false,
		},
		{
			name: "Empty clinic name",
			clinic: &clinic.Clinic{
				Name:        "",
				Address:     "123 Dental St",
				PhoneNumber: "1234567890",
				Email:       "contact@healthysmiles.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid clinic email",
			clinic: &clinic.Clinic{
				Name:        "Healthy Smiles",
				Address:     "123 Dental St",
				PhoneNumber: "1234567890",
				Email:       "invalid-email",
			},
			wantErr: true,
		},
		{
			name: "Short clinic address",
			clinic: &clinic.Clinic{
				Name:        "Healthy Smiles",
				Address:     "1",
				PhoneNumber: "1234567890",
				Email:       "contact@healthysmiles.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid clinic phone number",
			clinic: &clinic.Clinic{
				Name:        "Healthy Smiles",
				Address:     "123 Dental St",
				PhoneNumber: "12345",
				Email:       "contact@healthysmiles.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ClinicValidation(tt.clinic)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClinicValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
