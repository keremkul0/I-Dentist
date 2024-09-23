package tests

import (
	"bytes"
	"dental-clinic-system/api/signupClinic"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockSignUpClinicService struct {
	mock.Mock
}

func (m *MockSignUpClinicService) SignUpClinic(clinic models.Clinic, user models.User) (models.Clinic, error) {
	args := m.Called(clinic, user)
	return args.Get(0).(models.Clinic), args.Error(1)
}

func TestSignUpClinic_Success(t *testing.T) {
	mockService := new(MockSignUpClinicService)
	handler := signupClinic.NewSignUpClinicHandler(mockService)

	user := models.User{
		NationalID:    "123456789",
		Password:      "password",
		ClinicID:      1,
		Email:         "testuser@example.com",
		FirstName:     "Test",
		LastName:      "User",
		LastLogin:     time.Now(),
		IsActive:      true,
		EmailVerified: true,
		PhoneNumber:   "123-456-7890",
		PhoneVerified: true,
	}
	clinic := models.Clinic{
		Name:        "Test Clinic",
		Address:     "123 Test St",
		PhoneNumber: "123-456-7890",
	}
	signupData := map[string]interface{}{
		"user":   user,
		"clinic": clinic,
	}
	body, _ := json.Marshal(signupData)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	mockService.On("SignUpClinic", clinic, mock.AnythingOfType("models.User")).Return(clinic, nil)

	handler.SignUpClinic(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, "User, Group, and Clinic created successfully!", response["message"])
	assert.Equal(t, float64(clinic.ID), response["clinic_id"])
}

func TestSignUpClinic_InvalidJSON(t *testing.T) {
	mockService := new(MockSignUpClinicService)
	handler := signupClinic.NewSignUpClinicHandler(mockService)

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte("{invalid json}")))
	rr := httptest.NewRecorder()

	handler.SignUpClinic(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSignUpClinic_ServiceError(t *testing.T) {
	mockService := new(MockSignUpClinicService)
	handler := signupClinic.NewSignUpClinicHandler(mockService)

	user := models.User{
		NationalID:    "123456789",
		Password:      "password",
		Email:         "testuser@example.com",
		FirstName:     "Test",
		LastName:      "User",
		IsActive:      true,
		EmailVerified: true,
		PhoneNumber:   "123-456-7890",
		PhoneVerified: true,
	}
	clinic := models.Clinic{
		Name:        "Test Clinic",
		Address:     "123 Test St",
		PhoneNumber: "123-456-7890",
	}
	signupData := map[string]interface{}{
		"user":   user,
		"clinic": clinic,
	}
	body, _ := json.Marshal(signupData)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	mockService.On("SignUpClinic", clinic, mock.AnythingOfType("models.User")).Return(models.Clinic{}, assert.AnError)

	handler.SignUpClinic(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
