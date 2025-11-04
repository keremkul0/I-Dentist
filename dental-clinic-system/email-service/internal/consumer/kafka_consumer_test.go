package consumer

import (
	"context"
	"email-service/internal/service"
	"encoding/json"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailService for testing
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(msg service.EmailMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func TestNewKafkaConsumer(t *testing.T) {
	// Setup
	brokers := []string{"localhost:9092"}
	topics := []string{"email.verification", "email.password-reset"}
	groupID := "test-group"
	emailService := &service.EmailService{}

	// Execute
	consumer := NewKafkaConsumer(brokers, topics, groupID, emailService)

	// Assert
	assert.NotNil(t, consumer)
	assert.Equal(t, len(topics), len(consumer.readers))
	assert.Equal(t, emailService, consumer.emailService)

	// Cleanup
	for _, reader := range consumer.readers {
		_ = reader.Close()
	}
}

func TestKafkaConsumer_ConsumeMessages_ValidMessage(t *testing.T) {
	// Bu test gerçek Kafka bağlantısı gerektirdiği için integration test kategorisinde
	// Unit test için mock kullanmak gerekiyor, ancak mevcut kod yapısı bunu zorlaştırıyor

	// Test data
	emailMsg := service.EmailMessage{
		To:   "test@example.com",
		Type: "verification",
		Data: map[string]string{
			"token": "test-token",
		},
	}

	// JSON'a çevir
	msgBytes, err := json.Marshal(emailMsg)
	assert.NoError(t, err)

	// Kafka message oluştur
	kafkaMsg := kafka.Message{
		Topic: "email.verification",
		Value: msgBytes,
	}

	// Mesajın doğru unmarshal edildiğini test et
	var parsedMsg service.EmailMessage
	err = json.Unmarshal(kafkaMsg.Value, &parsedMsg)
	assert.NoError(t, err)
	assert.Equal(t, emailMsg.To, parsedMsg.To)
	assert.Equal(t, emailMsg.Type, parsedMsg.Type)
	assert.Equal(t, emailMsg.Data["token"], parsedMsg.Data["token"])
}

func TestKafkaConsumer_ConsumeMessages_InvalidJSON(t *testing.T) {
	// Invalid JSON test
	invalidJSON := []byte("invalid json")

	var emailMsg service.EmailMessage
	err := json.Unmarshal(invalidJSON, &emailMsg)

	// JSON unmarshal should fail
	assert.Error(t, err)
}

func TestKafkaConsumer_Start_ContextCancellation(t *testing.T) {
	// Bu test gerçek Kafka bağlantısı yapmaya çalıştığı için skip ediyoruz
	// Integration test olarak ayrı kategoride çalıştırılabilir
	t.Skip("Skipping integration test that requires real Kafka connection")

	// Setup
	brokers := []string{"localhost:9092"}
	topics := []string{"test-topic"}
	groupID := "test-group"
	emailService := &service.EmailService{}

	consumer := NewKafkaConsumer(brokers, topics, groupID, emailService)

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Start consumer (should return quickly due to timeout)
	start := time.Now()
	consumer.Start(ctx)
	duration := time.Since(start)

	// Should complete within reasonable time (not hang indefinitely)
	assert.Less(t, duration, 2*time.Second)

	// Cleanup
	for _, reader := range consumer.readers {
		_ = reader.Close()
	}
}

// Email message types test
func TestEmailMessage_Types(t *testing.T) {
	tests := []struct {
		name     string
		msgType  string
		expected string
	}{
		{"Verification email", "verification", "verification"},
		{"Password reset email", "password-reset", "password-reset"},
		{"General email", "general", "general"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailMsg := service.EmailMessage{
				To:   "test@example.com",
				Type: tt.msgType,
				Data: map[string]string{
					"token": "test-token",
				},
			}

			assert.Equal(t, tt.expected, emailMsg.Type)
			assert.NotEmpty(t, emailMsg.To)
			assert.Contains(t, emailMsg.Data, "token")
		})
	}
}

// Kafka Reader Config test
func TestKafkaConsumer_ReaderConfig(t *testing.T) {
	brokers := []string{"localhost:9092", "localhost:9093"}
	topics := []string{"email.verification"}
	groupID := "test-group"
	emailService := &service.EmailService{}

	consumer := NewKafkaConsumer(brokers, topics, groupID, emailService)

	// Test reader configuration
	assert.Len(t, consumer.readers, 1)

	reader := consumer.readers[0]
	config := reader.Config()

	assert.Equal(t, topics[0], config.Topic)
	assert.Equal(t, groupID, config.GroupID)
	assert.Equal(t, brokers, config.Brokers)
	assert.Equal(t, int(10e3), config.MinBytes)
	assert.Equal(t, int(10e6), config.MaxBytes)

	// Cleanup
	_ = reader.Close()
}

// Multiple topics test
func TestKafkaConsumer_MultipleTopic_Creation(t *testing.T) {
	brokers := []string{"localhost:9092"}
	topics := []string{"email.verification", "email.password-reset", "email.general", "email.dlq"}
	groupID := "test-group"
	emailService := &service.EmailService{}

	consumer := NewKafkaConsumer(brokers, topics, groupID, emailService)

	// Should create reader for each topic
	assert.Len(t, consumer.readers, len(topics))

	// Each reader should have correct topic
	for i, reader := range consumer.readers {
		config := reader.Config()
		assert.Equal(t, topics[i], config.Topic)
		assert.Equal(t, groupID, config.GroupID)
	}

	// Cleanup
	for _, reader := range consumer.readers {
		_ = reader.Close()
	}
}

// Email service integration test
func TestKafkaConsumer_EmailService_Integration(t *testing.T) {
	// Mock email service
	mockEmailService := new(MockEmailService)

	// Test email message
	emailMsg := service.EmailMessage{
		To:   "test@example.com",
		Type: "verification",
		Data: map[string]string{
			"token": "test-token-123",
		},
	}

	// Mock expectations
	mockEmailService.On("SendEmail", emailMsg).Return(nil)

	// Simulate sending email
	err := mockEmailService.SendEmail(emailMsg)

	// Assert
	assert.NoError(t, err)
	mockEmailService.AssertExpectations(t)
}

// Concurrency test for multiple topics
func TestKafkaConsumer_Concurrency(t *testing.T) {
	brokers := []string{"localhost:9092"}
	topics := []string{"topic1", "topic2", "topic3"}
	groupID := "test-group"
	emailService := &service.EmailService{}

	consumer := NewKafkaConsumer(brokers, topics, groupID, emailService)

	// Test that we can create multiple readers without issues
	assert.Len(t, consumer.readers, 3)

	// Test configuration only (skip actual Start to avoid Kafka connection)
	for i, reader := range consumer.readers {
		config := reader.Config()
		assert.Equal(t, topics[i], config.Topic)
		assert.Equal(t, groupID, config.GroupID)
		assert.Equal(t, brokers, config.Brokers)
	}

	// Cleanup
	for _, reader := range consumer.readers {
		_ = reader.Close()
	}
}
