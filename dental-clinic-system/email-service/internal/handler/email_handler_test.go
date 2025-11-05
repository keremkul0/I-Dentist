package handler

import (
	"bytes"
	"email-service/internal/service"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailService is a mock implementation of EmailServiceInterface
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(msg service.EmailMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func TestEmailHandler_SendEmail_Success(t *testing.T) {
	// Setup
	mockService := new(MockEmailService)
	handler := NewEmailHandler(mockService)

	app := fiber.New()
	app.Post("/send-email", handler.SendEmail)

	// Test data
	emailMsg := service.EmailMessage{
		To:   "test@example.com",
		Type: "verification",
		Data: map[string]string{
			"token": "test-token-123",
		},
	}

	// Mock expectations
	mockService.On("SendEmail", emailMsg).Return(nil)

	// Create request body
	reqBody, _ := json.Marshal(emailMsg)

	// Create test request
	req := httptest.NewRequest("POST", "/send-email", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check response body
	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Email sent successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestEmailHandler_SendEmail_InvalidJSON(t *testing.T) {
	// Setup
	mockService := new(MockEmailService)
	handler := NewEmailHandler(mockService)

	app := fiber.New()
	app.Post("/send-email", handler.SendEmail)

	// Create test request with invalid JSON
	req := httptest.NewRequest("POST", "/send-email", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Check response body
	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid request body", response["error"])
}

func TestEmailHandler_SendEmail_EmptyType(t *testing.T) {
	// Setup
	mockService := new(MockEmailService)
	handler := NewEmailHandler(mockService)

	app := fiber.New()
	app.Post("/send-email", handler.SendEmail)

	// Test data with empty type
	emailMsg := service.EmailMessage{
		To:   "test@example.com",
		Type: "", // Empty type
		Data: map[string]string{
			"token": "test-token-123",
		},
	}

	// Create request body
	reqBody, _ := json.Marshal(emailMsg)

	// Create test request
	req := httptest.NewRequest("POST", "/send-email", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Check response body
	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "Email type is required", response["error"])
}

func TestEmailHandler_SendEmail_ServiceError(t *testing.T) {
	// Setup
	mockService := new(MockEmailService)
	handler := NewEmailHandler(mockService)

	app := fiber.New()
	app.Post("/send-email", handler.SendEmail)

	// Test data
	emailMsg := service.EmailMessage{
		To:   "test@example.com",
		Type: "verification",
		Data: map[string]string{
			"token": "test-token-123",
		},
	}

	// Mock expectations - return error
	mockService.On("SendEmail", emailMsg).Return(errors.New("SMTP connection failed"))

	// Create request body
	reqBody, _ := json.Marshal(emailMsg)

	// Create test request
	req := httptest.NewRequest("POST", "/send-email", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Check response body
	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "Failed to send email", response["error"])
	assert.Equal(t, "SMTP connection failed", response["details"])

	mockService.AssertExpectations(t)
}

func TestEmailHandler_SendEmail_MissingToField(t *testing.T) {
	// Setup
	mockService := new(MockEmailService)
	handler := NewEmailHandler(mockService)

	app := fiber.New()
	app.Post("/send-email", handler.SendEmail)

	// Test data with missing To field
	emailMsg := service.EmailMessage{
		To:   "", // Missing To field
		Type: "verification",
		Data: map[string]string{
			"token": "test-token-123",
		},
	}

	// Mock expectations - service should return error for empty email
	mockService.On("SendEmail", emailMsg).Return(errors.New("email address is required"))

	// Create request body
	reqBody, _ := json.Marshal(emailMsg)

	// Create test request
	req := httptest.NewRequest("POST", "/send-email", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Check response body
	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, "Failed to send email", response["error"])
	assert.Contains(t, response["details"], "email address is required")

	mockService.AssertExpectations(t)
}

func TestEmailHandler_SendEmail_PasswordResetType(t *testing.T) {
	// Setup
	mockService := new(MockEmailService)
	handler := NewEmailHandler(mockService)

	app := fiber.New()
	app.Post("/send-email", handler.SendEmail)

	// Test data for password reset
	emailMsg := service.EmailMessage{
		To:   "test@example.com",
		Type: "password-reset",
		Data: map[string]string{
			"token": "reset-token-456",
		},
	}

	// Mock expectations
	mockService.On("SendEmail", emailMsg).Return(nil)

	// Create request body
	reqBody, _ := json.Marshal(emailMsg)

	// Create test request
	req := httptest.NewRequest("POST", "/send-email", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check response body
	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Email sent successfully", response["message"])

	mockService.AssertExpectations(t)
}
