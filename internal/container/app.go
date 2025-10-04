package container

import (
	"github.com/sony/sonyflake/v2"
	"github.com/tienhai2808/ecom_go/internal/config"
	customSf "github.com/tienhai2808/ecom_go/internal/snowflake"

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
	ProfileModule *ProfileModule
}

func NewContainer(db *gorm.DB, rdb *redis.Client, cfg *config.Config, rabbitChan *amqp091.Channel, sf *sonyflake.Sonyflake) *Container {
	sfg := customSf.NewSnowflakeGenerator(sf)
	userModule := NewUserContainer(db, sfg)
	authModule := NewAuthContainer(rdb, cfg, db, rabbitChan, sfg)
	addressModule := NewAddressContainer(db, sfg)
	productModule := NewProductContainer(db, sfg)
	imageModule := NewImageContainer(db, cfg)
	profileModule := NewProfileContainer(db)

	return &Container{
		userModule,
		authModule,
		addressModule,
		productModule,
		imageModule,
		profileModule,
	}
}
