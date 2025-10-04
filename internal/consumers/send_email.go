package consumers

import (
	"encoding/json"
	"fmt"
	"github.com/tienhai2808/ecom_go/internal/common"
	"github.com/tienhai2808/ecom_go/internal/types"
	"github.com/tienhai2808/ecom_go/internal/initialization"
	"github.com/tienhai2808/ecom_go/internal/rabbitmq"
	"github.com/tienhai2808/ecom_go/internal/smtp"
	"log"
)

func StartSendEmailConsumer(mqc *initialization.RabbitMQConn, mailer smtp.SMTPService) {
	if err := rabbitmq.ConsumeMessage(mqc.Chann, common.QueueName, common.Exchange, common.RoutingKey, func(body []byte) error {
		var emailMsg types.EmailMessage
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
