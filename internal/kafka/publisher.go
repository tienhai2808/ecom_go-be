package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

func PublishMessage(w *kafka.Writer, key, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return w.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

func PublishMessages(w *kafka.Writer, messages []kafka.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return w.WriteMessages(ctx, messages...)
}
