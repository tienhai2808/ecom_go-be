package initialization

import (
	"backend/config"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConn struct {
	Conn  *amqp091.Connection
	Chann *amqp091.Channel
}

func InitRabbitMQ(cfg *config.Config) (*RabbitMQConn, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Pass,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	fmt.Println(dsn)

	conn, err := amqp091.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("kết nối RabbitMQ thất bại: %w", err)
	}

	chann, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("mở channel thất bại: %w", err)
	}

	return &RabbitMQConn{
		conn,
		chann,
	}, nil
}

func (mq *RabbitMQConn) Close() {
	_ = mq.Chann.Close()
	_ = mq.Conn.Close()
}
