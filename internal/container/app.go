package container

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake/v2"
	customCld "github.com/tienhai2808/ecom_go/internal/cloudinary"
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/smtp"
	customSf "github.com/tienhai2808/ecom_go/internal/snowflake"
	"gorm.io/gorm"
)

type Container struct {
	UserModule     *UserModule
	AuthModule     *AuthModule
	AddressModule  *AddressModule
	ProductModule  *ProductModule
	ProfileModule  *ProfileModule
	CategoryModule *CategoryModule
	CartModule     *CartModule
	SMTPSvc        smtp.SMTPService
	CloudinarySvc  customCld.CloudinaryService
}

func NewContainer(db *gorm.DB, rdb *redis.Client, cfg *config.Config, rabbitChan *amqp091.Channel, sf *sonyflake.Sonyflake, cld *cloudinary.Cloudinary, es *elasticsearch.TypedClient) *Container {
	cSfg := customSf.NewSnowflakeGenerator(sf)
	smtp := smtp.NewSMTPService(cfg)
	cCld := customCld.NewCloudinaryService(cld)
	userModule := NewUserContainer(db, cSfg)
	authModule := NewAuthContainer(rdb, cfg, db, rabbitChan, cSfg)
	addressModule := NewAddressContainer(db, cSfg)
	productModule := NewProductContainer(db, rabbitChan, cSfg, es)
	profileModule := NewProfileContainer(db)
	categoryModule := NewCategoryContainer(db, cSfg)
	cartModule := NewCartModule(db, cSfg, es, cfg, rdb)

	return &Container{
		userModule,
		authModule,
		addressModule,
		productModule,
		profileModule,
		categoryModule,
		cartModule,
		smtp,
		cCld,
	}
}
