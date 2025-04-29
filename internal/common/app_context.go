package common

import (
	"backend/internal/config"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppContext struct {
	DB          *gorm.DB
	Redis       *redis.Client
	RabbitConn  *amqp091.Connection
	RabbitChan  *amqp091.Channel
	Config      *config.AppConfig
}