package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func PublishMessage(ch *amqp091.Channel, exchange, routingKey string, body []byte) error {
	_, err := ch.QueueDeclare(
		routingKey, 
		true,       
		false,      
		false,      
		false,      
		nil,        
	)
	if err != nil {
		return fmt.Errorf("không thể khai báo queue %s: %v", routingKey, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(
		ctx,
		exchange,  
		routingKey, 
		false,     
		false,     
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent, 
			Body:         body,
		},
	)
}