package consumer

import (
	"context"
	"email-service/internal/service"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	readers      []*kafka.Reader
	emailService *service.EmailService
}

func NewKafkaConsumer(brokers []string, topics []string, groupID string, emailService *service.EmailService) *KafkaConsumer {
	var readers []*kafka.Reader

	for _, topic := range topics {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
		readers = append(readers, reader)
	}

	return &KafkaConsumer{
		readers:      readers,
		emailService: emailService,
	}
}

func (kc *KafkaConsumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started...")

	var wg sync.WaitGroup

	// Her topic için ayrı goroutine
	for _, reader := range kc.readers {
		wg.Add(1)
		go func(r *kafka.Reader) {
			defer wg.Done()
			kc.consumeMessages(ctx, r)
		}(reader)
	}

	wg.Wait()
	log.Println("All consumers stopped")
}

func (kc *KafkaConsumer) consumeMessages(ctx context.Context, reader *kafka.Reader) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping consumer for topic: %s", reader.Config().Topic)
			if err := reader.Close(); err != nil {
				log.Printf("Error closing reader for topic %s: %v", reader.Config().Topic, err)
			}
			return
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from topic %s: %v", reader.Config().Topic, err)
				time.Sleep(1 * time.Second)
				continue
			}

			var emailMsg service.EmailMessage
			if err := json.Unmarshal(msg.Value, &emailMsg); err != nil {
				log.Printf("Error unmarshaling message from topic %s: %v", reader.Config().Topic, err)
				if commitErr := reader.CommitMessages(ctx, msg); commitErr != nil {
					log.Printf("Error committing failed message: %v", commitErr)
				}
				continue
			}

			// Email gönderme işlemi
			if err := kc.emailService.SendEmail(emailMsg); err != nil {
				log.Printf("Error sending email to %s: %v", emailMsg.To, err)
				// DLQ'ya gönder (opsiyonel)
			} else {
				log.Printf("Email sent successfully to: %s (type: %s)", emailMsg.To, emailMsg.Type)
			}

			// Message'ı commit et
			if err := reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}
