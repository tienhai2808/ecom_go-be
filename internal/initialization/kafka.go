package initialization

import (
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/tienhai2808/ecom_go/internal/config"
)

type KafkaClient struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

func InitKafka(cfg *config.Config) *KafkaClient {
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        "test-topic",
		Balancer:     &kafka.LeastBytes{}, // Load balancing strategy
		RequiredAcks: kafka.RequireAll,    // Đảm bảo message được ghi vào tất cả replica
		Compression:  kafka.Snappy,        // Nén message
		MaxAttempts:  5,                   // Số lần retry
		BatchSize:    100,                 // Batch messages
		BatchTimeout: 10 * time.Millisecond,
		Async:        false,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           cfg.Kafka.Brokers,
		Topic:             "test-topic",
		GroupID:           "test-consumer-group",
		MinBytes:          10e3,            // 10KB
		MaxBytes:          10e6,            // 10MB
		MaxWait:           1 * time.Second, // Max wait time
		ReadBackoffMin:    100 * time.Millisecond,
		ReadBackoffMax:    1 * time.Second,
		CommitInterval:    1 * time.Second, // Auto-commit interval
		StartOffset:       kafka.FirstOffset,
		SessionTimeout:    10 * time.Second,
		HeartbeatInterval: 3 * time.Second,
	})

	return &KafkaClient{
		w,
		r,
	}
}

func (k *KafkaClient) Close() {
	_ = k.Writer.Close()
	_ = k.Reader.Close()
}
