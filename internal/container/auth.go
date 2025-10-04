package container

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/smtp"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthModule struct {
	AuthHandler *handler.AuthHandler
	SMTPService smtp.SMTPService
}

func NewAuthContainer(rdb *redis.Client, cfg *config.Config, db *gorm.DB, rabbitChan *amqp091.Channel, sfg snowflake.SnowflakeGenerator) *AuthModule {
	mailer := smtp.NewSMTPService(cfg)
	authRepo := repoImpl.NewAuthRepository(rdb, cfg)
	userRepo := repoImpl.NewUserRepository(db)
	profileRepo := repoImpl.NewProfileRepository(db)
	authSvc := serviceImpl.NewAuthService(userRepo, authRepo, profileRepo, rabbitChan, cfg, sfg)
	userSvc := serviceImpl.NewUserService(userRepo, profileRepo, sfg)
	authHandler := handler.NewAuthHandler(authSvc, userSvc, cfg)

	return &AuthModule{
		authHandler,
		mailer,
	}
}
