package mq

import "github.com/rabbitmq/amqp091-go"

func ConsumeMessage(ch *amqp091.Channel, queueName, exchange, routingKey string, handler func([]byte) error) error {
	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return err
	}

	if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		return err
	}

	if err := ch.QueueBind(queueName, routingKey, exchange, false, nil); err != nil {
		return err
	}

	if err := ch.Qos(5, 0, false); err != nil {
		return err
	}

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		go func(workerID int) {
			for msg := range msgs {
				processWithRetry(msg.Body, handler, workerID)
			}
		}(i)
	}

	return nil
}

