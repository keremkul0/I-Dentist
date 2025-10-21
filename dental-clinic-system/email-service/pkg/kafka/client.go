package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

func NewClient(config Config) *Client {
	// Producer (Writer) konfigürasyonu
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		BatchTimeout: 10 * time.Millisecond,
	}

	// Consumer (Reader) konfigürasyonu
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.Brokers,
		Topic:    config.Topic,
		GroupID:  config.GroupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  1 * time.Second,
	})

	return &Client{
		writer: writer,
		reader: reader,
	}
}

// Message producer - Ana backend'den email servisine mesaj gönderir
func (c *Client) ProduceMessage(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	}

	return c.writer.WriteMessages(ctx, message)
}

// Message consumer - Email servisinde mesajları dinler
func (c *Client) ConsumeMessages(ctx context.Context, handler func([]byte) error) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopping...")
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error fetching message: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// Mesajı işle
			if err := handler(msg.Value); err != nil {
				log.Printf("Error handling message: %v", err)
			}

			// Mesajı commit et (başarılı işlendikten sonra)
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

// Topic oluşturma - İlk kurulum için
func (c *Client) CreateTopic(ctx context.Context, topic string, partitions int, replicationFactor int) error {
	conn, err := kafka.DialLeader(ctx, "tcp", c.writer.Addr.String(), topic, 0)
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	return conn.CreateTopics(topicConfig)
}

// Bağlantıları kapat
func (c *Client) Close() error {
	if err := c.writer.Close(); err != nil {
		log.Printf("Error closing writer: %v", err)
		return err
	}

	if err := c.reader.Close(); err != nil {
		log.Printf("Error closing reader: %v", err)
		return err
	}

	return nil
}

// Health check - Kafka bağlantısını kontrol et
func (c *Client) HealthCheck(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", c.writer.Addr.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
