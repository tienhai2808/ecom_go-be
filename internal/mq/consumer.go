// In mq/consumer.go
package mq

import (
	"backend/internal/common"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func ConsumeMessages(ch *amqp091.Channel, queueName string) (<-chan amqp091.Delivery, error) {
	// Ensure the queue exists
	_, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("không thể khai báo queue %s: %v", queueName, err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack (set to false to manually acknowledge)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("không thể consume từ queue %s: %v", queueName, err)
	}
	return msgs, nil
}

func StartEmailConsumer(ch *amqp091.Channel, emailSender common.EmailSender) {
	queueName := "email_queue"

	// Ensure the queue exists
	_, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("❤️ Không thể khai báo queue %s: %v", queueName, err)
	}

	msgs, err := ConsumeMessages(ch, queueName)
	if err != nil {
		log.Fatalf("❤️ Không thể consume từ queue %s: %v", queueName, err)
	}

	log.Printf("💚 Bắt đầu lắng nghe email queue: %s", queueName)

	go func() {
		for msg := range msgs {
			log.Printf("📥 Nhận được email message")

			var emailMsg EmailMessage
			if err := json.Unmarshal(msg.Body, &emailMsg); err != nil {
				log.Printf("❤️ Lỗi parse email message: %v", err)
				// Give some time before requeueing to avoid rapid requeuing of malformed messages
				time.Sleep(1 * time.Second)
				msg.Nack(false, true) // Nack and requeue
				continue
			}

			// Validate email message fields
			if emailMsg.To == "" || emailMsg.Subject == "" {
				log.Printf("❤️ Email message thiếu thông tin cần thiết: To=%s, Subject=%s",
					emailMsg.To, emailMsg.Subject)
				// This is likely a permanent error, don't requeue
				msg.Nack(false, false)
				continue
			}

			log.Printf("📧 Gửi email đến: %s, chủ đề: %s", emailMsg.To, emailMsg.Subject)

			// Retry mechanism for sending emails
			maxRetries := 3
			var sendErr error

			for i := 0; i < maxRetries; i++ {
				sendErr = emailSender.SendEmail(emailMsg.To, emailMsg.Subject, emailMsg.Body)
				if sendErr == nil {
					break
				}

				log.Printf("❤️ Lần thử %d/%d: Lỗi gửi email: %v", i+1, maxRetries, sendErr)
				if i < maxRetries-1 {
					time.Sleep(2 * time.Second)
				}
			}

			if sendErr != nil {
				log.Printf("❤️ Đã thử %d lần nhưng không thể gửi email: %v", maxRetries, sendErr)
				// Nack the message and requeue it to try again later
				msg.Nack(false, true)
				continue
			}

			log.Printf("💚 Đã gửi email thành công đến: %s", emailMsg.To)
			msg.Ack(false)
		}

		log.Println("❤️ Email consumer channel closed. Waiting for reconnection...")
	}()
}
