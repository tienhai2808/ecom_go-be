package mq

import (
	"backend/internal/config"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ(cfg *config.AppConfig) (*amqp091.Connection, *amqp091.Channel, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Pass,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("❤️ Không thể kết nối đến RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close() 
		return nil, nil, fmt.Errorf("❤️ Không thể mở channel RabbitMQ: %v", err)
	}

	return conn, channel, nil
}

func CloseRabbitMQ(conn *amqp091.Connection, channel *amqp091.Channel) {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
