package container

import (
	"backend/internal/config"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	AuthModule    *AuthModule
	ProductModule *ProductModule
}

func NewContainer(db *gorm.DB, redis *redis.Client, config *config.AppConfig, rabbitChan *amqp091.Channel) (*Container, error) {
	authModule := NewAuthContainer(redis, config, db, rabbitChan)
	productModule := NewProductContainer(db)

	return &Container{
		AuthModule:    authModule,
		ProductModule: productModule,
	}, nil
}