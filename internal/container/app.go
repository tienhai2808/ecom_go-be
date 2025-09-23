package container

import (
	"github.com/tienhai2808/ecom_go/config"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	UserModule    *UserModule
	AuthModule    *AuthModule
	AddressModule *AddressModule
	ProductModule *ProductModule
	ImageModule   *ImageModule
}

func NewContainer(db *gorm.DB, redis *redis.Client, config *config.Config, rabbitChan *amqp091.Channel) *Container {
	userModule := NewUserContainer(db)
	authModule := NewAuthContainer(redis, config, db, rabbitChan)
	addressModule := NewAddressContainer(db)
	productModule := NewProductContainer(db)
	imageModule := NewImageContainer(db, config)

	return &Container{
		UserModule:    userModule,
		AuthModule:    authModule,
		AddressModule: addressModule,
		ProductModule: productModule,
		ImageModule:   imageModule,
	}
}
