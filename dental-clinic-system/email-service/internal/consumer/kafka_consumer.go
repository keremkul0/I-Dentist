package consumer

import (
	"context"
	"email-service/internal/service"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type EmailMessage struct {
	To   string
	Type string
}

type KafkaConsumer struct {
	reader       *kafka.Reader
	emailService *service.EmailService
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, emailService *service.EmailService) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
		//MinBytes: 10e3, // 10KB
		//MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		reader:       reader,
		emailService: emailService,
	}
}

func (kc *KafkaConsumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer...")
			if err := kc.reader.Close(); err != nil {
				log.Printf("Error closing reader: %v", err)
			}
			return
		default:
			msg, err := kc.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			var emailMsg service.EmailMessage
			if err := json.Unmarshal(msg.Value, &emailMsg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				if commitErr := kc.reader.CommitMessages(ctx, msg); commitErr != nil {
					log.Printf("Error committing failed message: %v", commitErr)
				}
				continue
			}

			// Email gönderme işlemi
			if err := kc.emailService.SendEmail(emailMsg); err != nil {
				log.Printf("Error sending email to %s: %v", emailMsg.To, err)
			} else {
				log.Printf("Email sent successfully to: %s (type: %s)", emailMsg.To, emailMsg.Type)
			}

			// Message'ı commit et
			if err := kc.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}
