package server

import (
	"backend/internal/cache"
	"backend/internal/config"
	"backend/internal/container"
	"backend/internal/database"
	"backend/internal/mq"
	"backend/internal/router"
	"backend/internal/smtp"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Config     *config.AppConfig
	DB         *gorm.DB
	Redis      *redis.Client
	RabbitConn *amqp091.Connection
	RabbitChan *amqp091.Channel
	Container  *container.Container
	Router     *gin.Engine
}

func NewApplication() *Application {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("❤️ Lỗi khi load cấu hình app: %v", err)
	}

	db, err := database.ConnectToDatabase(cfg)
	if err != nil {
		log.Fatalf("❤️ Lỗi kết nối tới database: %v", err)
	}

	redisClient, err := cache.ConnectToRedis(cfg)
	if err != nil {
		log.Fatalf("❤️ Lỗi kết nối tới redis: %v", err)
	}

	rabbitConn, rabbitChan, err := mq.ConnectToRabbitMQ(cfg)
	if err != nil {
		log.Fatalf("❤️ Lỗi kết nối tới rabbitmq: %v", err)
	}

	emailSender := smtp.NewSMTPSender(cfg)
	mq.StartEmailConsumer(rabbitChan, emailSender)

	con, err := container.NewContainer(db, redisClient, cfg, rabbitChan)
	if err != nil {
		log.Fatalf("❤️ Lỗi tạo container: %v", err)
	}

	router := gin.Default()
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("❤️ Could not set trusted proxies: %v", err)
	}

	config.SetupCORS(router)

	app := &Application{
		Config:     cfg,
		DB:         db,
		Redis:      redisClient,
		RabbitConn: rabbitConn,
		RabbitChan: rabbitChan,
		Container:  con,
		Router:     router,
	}

	app.initRoutes()

	return app
}

func (app *Application) initRoutes() {
	api := app.Router.Group(app.Config.App.ApiPrefix)
	router.NewUserRouter(api, app.Config, app.Container.UserModule.UserRepository, app.Container.UserModule.UserHandler)
	router.NewAuthRouter(api, app.Config, app.Container.UserModule.UserRepository, app.Container.AuthModule.AuthHandler)
	router.NewAddressRouter(api, app.Config, app.Container.UserModule.UserRepository, app.Container.AddressModule.AddressHandler)
	router.NewProductRouter(api, app.Container.ProductModule.ProductHandler)
}

func (app *Application) Run() {
	fmt.Println("💚 Kết nối MySQL thành công")
	fmt.Println("💚 Kết nối Redis thành công")
	fmt.Println("💚 Kết nối RabbitMQ thành công")
	addr := app.Config.App.Host + ":" + app.Config.App.Port
	if err := app.Router.Run(addr); err != nil {
		log.Fatalf("❤️ Không thể khởi động server: %v", err)
	}
}

func (app *Application) Close() {
	if sqlDB, err := app.DB.DB(); err == nil {
		sqlDB.Close()
	}
	if app.Redis != nil {
		app.Redis.Close()
	}
	mq.CloseRabbitMQ(app.RabbitConn, app.RabbitChan)
}
