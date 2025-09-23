package consumers

import (
	"backend/internal/common"
	"backend/internal/dto"
	"backend/internal/initialization"
	"backend/internal/mq"
	"backend/internal/smtp"
	"encoding/json"
	"fmt"
	"log"
)

func StartSendEmailConsumer(mqc *initialization.RabbitMQConn, mailer smtp.SMTPService) {
	if err := mq.ConsumeMessage(mqc.Chann, common.QueueName, common.Exchange, common.RoutingKey, func(body []byte) error {
		var emailMsg dto.EmailMessage
		if err := json.Unmarshal(body, &emailMsg); err != nil {
			return fmt.Errorf("chuyển đổi tin nhắn email thất bại: %w", err)
		}

		if err := mailer.SendEmail(emailMsg.To, emailMsg.Subject, emailMsg.Body); err != nil {
			return fmt.Errorf("gửi email thất bại: %w", err)
		}

		log.Printf("Đã gửi email thành công tới: %s", emailMsg.To)
		return nil
	}); err != nil {
		log.Printf("Lỗi khởi tạo email consumer: %v", err)
	}
}