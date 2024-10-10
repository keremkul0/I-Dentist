package signupClinic_test

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

func TestSignUpClinic(t *testing.T) {
	mockService := new(MockSignUpClinicService)
	handler := signupClinic.NewSignUpClinicHandler(mockService)

	clinic := models.Clinic{
		Model:       gorm.Model{ID: 1},
		Name:        "Test Clinic",
		Address:     "123 Test St",
		PhoneNumber: "123-456-7890",
		Email:       "testclinic@example.com",
	}
	user := models.User{
		NationalID:    "12345678901",
		Password:      "hashedpassword",
		Email:         "user@example.com",
		EmailVerified: true,
		FirstName:     "John",
		LastName:      "Doe",
		LastLogin:     time.Now(),
		IsActive:      true,
		PhoneNumber:   "987-654-3210",
		PhoneVerified: true,
		Roles:         []*models.Role{{Name: "Admin"}, {Name: "User"}},
	}

	userGetModel := models.UserGetModel{
		NationalID:  "12345678901",
		Email:       "user@example.com",
		FirstName:   "John",
		LastName:    "Doe",
		LastLogin:   time.Now(),
		IsActive:    true,
		PhoneNumber: "987-654-3210",
		Roles:       []*models.Role{{Name: "Admin"}, {Name: "User"}},
	}

	mockService.On("SignUpClinic", clinic, user).Return(clinic, userGetModel, nil)

	body := map[string]interface{}{
		"clinic": clinic,
		"user":   user,
	}
	bodyBytes, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.SignUpClinic(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	// Check the response body
	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	// Assert clinic information in response
	clinicResp := responseBody["clinic"].(map[string]interface{})
	assert.Equal(t, clinic.Name, clinicResp["name"])
	assert.Equal(t, clinic.Address, clinicResp["address"])
	assert.Equal(t, clinic.PhoneNumber, clinicResp["phoneNumber"])
	assert.Equal(t, clinic.Email, clinicResp["email"])

	// Assert user information in response
	userResp := responseBody["user"].(map[string]interface{})
	assert.Equal(t, user.NationalID, userResp["nationalID"])
	assert.Equal(t, user.Email, userResp["email"])
	assert.Equal(t, user.FirstName, userResp["firstName"])
	assert.Equal(t, user.LastName, userResp["lastName"])
	assert.Equal(t, user.PhoneNumber, userResp["phoneNumber"])

	// Assert roles in response
	rolesResp := userResp["roles"].([]interface{})
	assert.Equal(t, "Admin", rolesResp[0].(map[string]interface{})["name"])
	assert.Equal(t, "User", rolesResp[1].(map[string]interface{})["name"])

	// Verify that the mock service's method was called
	mockService.AssertCalled(t, "SignUpClinic", clinic, user)
}
