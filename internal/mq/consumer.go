package mq

import (
	"backend/internal/smtp"
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
	_, err := ch.QueueDeclare(
		queueName, 
		true,      
		false,     
		false,     
		false,     
		nil,       
	)
	if err != nil {
		return nil, fmt.Errorf("không thể khai báo queue %s: %v", queueName, err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",   
		false, 
		false, 
		false, 
		false, 
		nil,   
	)
	if err != nil {
		return nil, fmt.Errorf("không thể consume từ queue %s: %v", queueName, err)
	}
	return msgs, nil
}

func StartEmailConsumer(ch *amqp091.Channel, emailSender smtp.EmailSender) {
	queueName := "email_queue"

	_, err := ch.QueueDeclare(
		queueName, 
		true,      
		false,     
		false,     
		false,     
		nil,       
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
				time.Sleep(1 * time.Second)
				msg.Nack(false, true) 
				continue
			}

			if emailMsg.To == "" || emailMsg.Subject == "" {
				log.Printf("❤️ Email message thiếu thông tin cần thiết: To=%s, Subject=%s",
					emailMsg.To, emailMsg.Subject)
				msg.Nack(false, false)
				continue
			}

			log.Printf("📧 Gửi email đến: %s, chủ đề: %s", emailMsg.To, emailMsg.Subject)

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
				msg.Nack(false, true)
				continue
			}

			log.Printf("💚 Đã gửi email thành công đến: %s", emailMsg.To)
			msg.Ack(false)
		}

		log.Println("❤️ Email consumer channel closed. Waiting for reconnection...")
	}()
}
