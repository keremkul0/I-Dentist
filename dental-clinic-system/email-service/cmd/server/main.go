package main

import (
	"context"
	"email-service/internal/consumer"
	"email-service/internal/service"
	"email-service/pkg/email"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// SMTP mailer oluştur
	mailer := email.NewSMTPMailer()

	// Email service oluştur
	emailService := service.NewEmailService(mailer)

	// Kafka consumer oluştur
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	if len(brokers) == 0 || topic == "" || groupID == "" {
		log.Fatal("Kafka configuration missing. Please set KAFKA_BROKERS, KAFKA_TOPIC, and KAFKA_GROUP_ID")
	}

	kafkaConsumer := consumer.NewKafkaConsumer(brokers, topic, groupID, emailService)

	// Context oluştur
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown için signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Kafka consumer'ı başlat
	go kafkaConsumer.Start(ctx)

	log.Println("Email service started successfully")
	log.Printf("Listening to Kafka topic: %s", topic)
	log.Printf("Kafka brokers: %v", brokers)

	// Shutdown signal'ını bekle
	<-sigChan
	log.Println("Shutdown signal received, stopping email service...")

	cancel()
	log.Println("Email service stopped")
}
