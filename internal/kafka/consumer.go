package kafka

import (
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

func ConsumeMessages(ctx context.Context, r *kafka.Reader, handler func(kafka.Message) error) {
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			if err == kafka.ErrGroupClosed || err.Error() == "EOF" {
				return
			}
			log.Printf("[KAFKA] Lỗi fetch: %v", err)
			continue
		}
		log.Printf("[KAFKA CONSUMER] ✅ Nhận message: key=%s, partition=%d, offset=%d",
			string(msg.Key),
			msg.Partition,
			msg.Offset,
		)

		if err := handler(msg); err != nil {
			log.Printf("Error processing message: %v", err)
			continue
		}

		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("Failed to commit message: %v", err)
		}
	}
}

func MessageHandler(msg kafka.Message) error {
	log.Printf("Received message: key=%s, value=%s, partition=%d, offset=%d",
		string(msg.Key),
		string(msg.Value),
		msg.Partition,
		msg.Offset,
	)
	return nil
}
