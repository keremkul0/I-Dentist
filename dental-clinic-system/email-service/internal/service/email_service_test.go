package service

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/gomail.v2"
)

// MockMailer implements the Mailer interface for testing
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendMail(message gomail.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

// setupTestTemplates creates test template files
func setupTestTemplates(t *testing.T) {
	err := os.MkdirAll("templates", 0755)
	assert.NoError(t, err)

	verificationTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>Email Verification</title>
</head>
<body>
    <h1>Verify Your Email</h1>
    <p>Click the link below to verify your email:</p>
    <a href="{{.VERIFY_LINK}}">Verify Email</a>
</body>
</html>`

	passwordResetTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>Password Reset</title>
</head>
<body>
    <h1>Reset Your Password</h1>
    <p>Click the link below to reset your password:</p>
    <a href="{{.RESET_LINK}}">Reset Password</a>
</body>
</html>`

	err = os.WriteFile("templates/verification_email.html", []byte(verificationTemplate), 0644)
	assert.NoError(t, err)

	err = os.WriteFile("templates/password_reset_email.html", []byte(passwordResetTemplate), 0644)
	assert.NoError(t, err)
}

// cleanupTestTemplates removes test template files
func cleanupTestTemplates() {
	os.RemoveAll("templates")
}

func TestNewEmailService(t *testing.T) {
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	assert.NotNil(t, service)
	assert.Equal(t, mockMailer, service.mailer)
}

func TestEmailService_SendEmail_VerificationType(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Setup test templates
	setupTestTemplates(t)
	defer cleanupTestTemplates()

	// Set required environment variables
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	os.Setenv("SMTP_FROM", "test@example.com")
	defer func() {
		os.Unsetenv("FRONTEND_URL")
		os.Unsetenv("SMTP_FROM")
	}()

	// Mock expectations
	mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(nil)

	// Test data
	emailMsg := EmailMessage{
		To:   "user@example.com",
		Type: "verification",
		Data: map[string]string{
			"token": "test-token-123",
		},
	}

	// Execute
	err := service.SendEmail(emailMsg)

	// Assert
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

func TestEmailService_SendEmail_PasswordResetType(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Setup test templates
	setupTestTemplates(t)
	defer cleanupTestTemplates()

	// Set required environment variables
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	os.Setenv("SMTP_FROM", "test@example.com")
	defer func() {
		os.Unsetenv("FRONTEND_URL")
		os.Unsetenv("SMTP_FROM")
	}()

	// Mock expectations
	mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(nil)

	// Test data
	emailMsg := EmailMessage{
		To:   "user@example.com",
		Type: "password_reset",
		Data: map[string]string{
			"token": "reset-token-456",
		},
	}

	// Execute
	err := service.SendEmail(emailMsg)

	// Assert
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

func TestEmailService_SendEmail_MailerError(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Setup test templates
	setupTestTemplates(t)
	defer cleanupTestTemplates()

	// Set required environment variables
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	os.Setenv("SMTP_FROM", "test@example.com")
	defer func() {
		os.Unsetenv("FRONTEND_URL")
		os.Unsetenv("SMTP_FROM")
	}()

	// Mock expectations - return error
	expectedError := errors.New("SMTP connection failed")
	mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(expectedError)

	// Test data
	emailMsg := EmailMessage{
		To:   "user@example.com",
		Type: "verification",
		Data: map[string]string{
			"token": "test-token-123",
		},
	}

	// Execute
	err := service.SendEmail(emailMsg)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockMailer.AssertExpectations(t)
}

func TestEmailService_SendEmail_MissingToken(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Setup test templates
	setupTestTemplates(t)
	defer cleanupTestTemplates()

	// Set required environment variables
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	os.Setenv("SMTP_FROM", "test@example.com")
	defer func() {
		os.Unsetenv("FRONTEND_URL")
		os.Unsetenv("SMTP_FROM")
	}()

	// Mock expectations
	mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(nil)

	// Test data - missing token
	emailMsg := EmailMessage{
		To:   "user@example.com",
		Type: "verification",
		Data: map[string]string{}, // Empty data
	}

	// Execute
	err := service.SendEmail(emailMsg)

	// Assert - should still work but with empty token
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

func TestEmailService_SendVerificationEmail(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Setup test templates
	setupTestTemplates(t)
	defer cleanupTestTemplates()

	// Set required environment variables
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	os.Setenv("SMTP_FROM", "test@example.com")
	defer func() {
		os.Unsetenv("FRONTEND_URL")
		os.Unsetenv("SMTP_FROM")
	}()

	// Mock expectations
	mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(nil)

	// Execute
	err := service.sendVerificationEmail("user@example.com", "test-token")

	// Assert
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

func TestEmailService_SendPasswordResetEmail(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Setup test templates
	setupTestTemplates(t)
	defer cleanupTestTemplates()

	// Set required environment variables
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	os.Setenv("SMTP_FROM", "test@example.com")
	defer func() {
		os.Unsetenv("FRONTEND_URL")
		os.Unsetenv("SMTP_FROM")
	}()

	// Mock expectations
	mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(nil)

	// Execute
	err := service.sendPasswordResetEmail("user@example.com", "reset-token")

	// Assert
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

func TestEmailService_SendTemplateEmail_TemplateNotFound(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Set required environment variables
	os.Setenv("SMTP_FROM", "test@example.com")
	defer os.Unsetenv("SMTP_FROM")

	// Execute with non-existent template
	err := service.sendTemplateEmail(
		"user@example.com",
		"Test Subject",
		"templates/non_existent_template.html",
		map[string]string{"KEY": "value"},
	)

	// Assert
	assert.Error(t, err)
	// Mailer should not be called since template parsing fails
	mockMailer.AssertNotCalled(t, "SendMail")
}

func TestEmailService_SendTemplateEmail_MissingEnvironmentVariable(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Setup test templates
	setupTestTemplates(t)
	defer cleanupTestTemplates()

	// Ensure SMTP_FROM is not set
	os.Unsetenv("SMTP_FROM")

	// Mock expectations - should still be called even with empty FROM
	mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(nil)

	// Execute
	err := service.sendTemplateEmail(
		"user@example.com",
		"Test Subject",
		"templates/verification_email.html",
		map[string]string{"VERIFY_LINK": "http://example.com/verify"},
	)

	// Assert - should work even with empty SMTP_FROM (gomail will handle it)
	assert.NoError(t, err)
	mockMailer.AssertExpectations(t)
}

func TestEmailService_SendTemplateEmail_InvalidTemplateData(t *testing.T) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Set required environment variables
	os.Setenv("SMTP_FROM", "test@example.com")
	defer os.Unsetenv("SMTP_FROM")

	// Create test templates directory
	err := os.MkdirAll("templates", 0755)
	assert.NoError(t, err)
	defer os.RemoveAll("templates")

	// Create a temporary template with invalid Go template syntax
	tempTemplate := `<!DOCTYPE html><html><body>{{.MissingClosingBrace</body></html>`
	tempFile := "templates/invalid_template.html"
	err = os.WriteFile(tempFile, []byte(tempTemplate), 0644)
	assert.NoError(t, err)

	// Execute
	err = service.sendTemplateEmail(
		"user@example.com",
		"Test Subject",
		tempFile,
		map[string]string{"KEY": "value"},
	)

	// Assert - template parsing should fail
	assert.Error(t, err)
	mockMailer.AssertNotCalled(t, "SendMail")
}

// Benchmark tests
func BenchmarkEmailService_SendEmail(b *testing.B) {
	// Setup
	mockMailer := &MockMailer{}
	service := NewEmailService(mockMailer)

	// Setup test templates
	err := os.MkdirAll("templates", 0755)
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll("templates")

	verificationTemplate := `<!DOCTYPE html><html><body><h1>Verify Your Email</h1><a href="{{.VERIFY_LINK}}">Verify Email</a></body></html>`
	err = os.WriteFile("templates/verification_email.html", []byte(verificationTemplate), 0644)
	if err != nil {
		b.Fatal(err)
	}

	// Set required environment variables
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	os.Setenv("SMTP_FROM", "test@example.com")
	defer func() {
		os.Unsetenv("FRONTEND_URL")
		os.Unsetenv("SMTP_FROM")
	}()

	// Mock expectations
	mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(nil)

	emailMsg := EmailMessage{
		To:   "user@example.com",
		Type: "verification",
		Data: map[string]string{
			"token": "test-token-123",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.SendEmail(emailMsg)
	}
}

// Table-driven tests for multiple email types
func TestEmailService_SendEmail_MultipleTypes(t *testing.T) {
	testCases := []struct {
		name      string
		emailType string
		token     string
		expectErr bool
	}{
		{
			name:      "Verification email",
			emailType: "verification",
			token:     "verify-token-123",
			expectErr: false,
		},
		{
			name:      "Password reset email",
			emailType: "password_reset",
			token:     "reset-token-456",
			expectErr: false,
		},
		{
			name:      "Unknown email type - defaults to password reset",
			emailType: "unknown_type",
			token:     "unknown-token-789",
			expectErr: false,
		},
		{
			name:      "Empty email type - defaults to password reset",
			emailType: "",
			token:     "empty-token-000",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockMailer := &MockMailer{}
			service := NewEmailService(mockMailer)

			// Setup test templates
			setupTestTemplates(t)
			defer cleanupTestTemplates()

			// Set required environment variables
			os.Setenv("FRONTEND_URL", "http://localhost:3000")
			os.Setenv("SMTP_FROM", "test@example.com")
			defer func() {
				os.Unsetenv("FRONTEND_URL")
				os.Unsetenv("SMTP_FROM")
			}()

			// Mock expectations
			mockMailer.On("SendMail", mock.AnythingOfType("gomail.Message")).Return(nil)

			// Test data
			emailMsg := EmailMessage{
				To:   "user@example.com",
				Type: tc.emailType,
				Data: map[string]string{
					"token": tc.token,
				},
			}

			// Execute
			err := service.SendEmail(emailMsg)

			// Assert
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockMailer.AssertExpectations(t)
			}
		})
	}
}
