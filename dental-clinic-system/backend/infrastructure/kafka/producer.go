package kafka

import (
	"context"
	"dental-clinic-system/infrastructure/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type EmailMessage struct {
	Type string            `json:"type"`
	To   string            `json:"to"`
	Data map[string]string `json:"data"`
}

type EmailProducer interface {
	SendVerificationEmail(email, token string) error
	SendPasswordResetEmail(email, token string) error
	Close() error
}

type kafkaEmailProducer struct {
	writer *kafka.Writer
	config *config.KafkaConfig
}

func NewEmailProducer(cfg *config.KafkaConfig) EmailProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		RequiredAcks: kafka.RequireOne,
		Compression:  kafka.Gzip,
	}

	return &kafkaEmailProducer{
		writer: writer,
		config: cfg,
	}
}

func (p *kafkaEmailProducer) SendVerificationEmail(email, token string) error {
	message := EmailMessage{
		Type: "verification",
		To:   email,
		Data: map[string]string{
			"token": token,
		},
	}

	return p.sendMessage(p.config.VerificationTopic, message)
}

func (p *kafkaEmailProducer) SendPasswordResetEmail(email, token string) error {
	message := EmailMessage{
		Type: "password-reset",
		To:   email,
		Data: map[string]string{
			"token": token,
		},
	}

	return p.sendMessage(p.config.PasswordResetTopic, message)
}

func (p *kafkaEmailProducer) sendMessage(topic string, message EmailMessage) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Str("topic", topic).Msg("Failed to marshal email message")
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMessage := kafka.Message{
		Topic: topic,
		Key:   []byte(message.To), // Use email as key for partitioning
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "type", Value: []byte(message.Type)},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = p.writer.WriteMessages(ctx, kafkaMessage)
	if err != nil {
		log.Error().Err(err).
			Str("topic", topic).
			Str("email", message.To).
			Str("type", message.Type).
			Msg("Failed to send email message to Kafka")
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Info().
		Str("topic", topic).
		Str("email", message.To).
		Str("type", message.Type).
		Msg("Email message sent to Kafka successfully")

	return nil
}

func (p *kafkaEmailProducer) Close() error {
	return p.writer.Close()
}
