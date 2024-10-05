package tests

import (
	"bytes"
	"dental-clinic-system/api/signupClinic"
	"dental-clinic-system/models"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service
type MockSignUpClinicService struct {
	mock.Mock
}

func (m *MockSignUpClinicService) SignUpClinic(clinic models.Clinic, user models.User) (models.Clinic, models.UserGetModel, error) {
	args := m.Called(clinic, user)
	return args.Get(0).(models.Clinic), args.Get(1).(models.UserGetModel), args.Error(2)
}

func TestSignUpClinicHandler_SignUpClinic(t *testing.T) {
	mockService := new(MockSignUpClinicService)
	handler := signupClinic.NewSignUpClinicHandler(mockService)

	clinic := models.Clinic{Model: gorm.Model{ID: 1}, Name: "Test Clinic", Address: "123 Main St", PhoneNumber: "123-456-7890", Email: "testclinic@example.com"}

	user := models.User{Model: gorm.Model{ID: 1}, NationalID: "123456789", Password: "hashedpassword", ClinicID: 1, Email: "testuser@example.com", EmailVerified: true, FirstName: "Test", LastName: "User", LastLogin: time.Now(), IsActive: true, PhoneNumber: "123-456-7890", PhoneVerified: true, Roles: []*models.Role{}}
	userGetModel := models.UserGetModel{Model: gorm.Model{ID: 1}, NationalID: "123456789", ClinicID: 1, Email: "testuser@example.com", FirstName: "Test", LastName: "User", LastLogin: time.Now(), IsActive: true, PhoneNumber: "123-456-7890", Roles: []*models.Role{}}
	mockService.On("SignUpClinic", clinic, user).Return(clinic, userGetModel, nil)

	body := struct {
		Clinic models.Clinic `json:"clinic"`
		User   models.User   `json:"user"`
	}{clinic, user}
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.SignUpClinic(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response struct {
		Clinic models.Clinic       `json:"clinic"`
		User   models.UserGetModel `json:"user"`
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, clinic, response.Clinic)
	assert.Equal(t, userGetModel, response.User)
}
