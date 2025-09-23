package container

import (
	"backend/config"
	"backend/internal/handler"
	repoImpl "backend/internal/repository/implement"
	serviceImpl "backend/internal/service/implement"
	"backend/internal/smtp"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthModule struct {
	AuthHandler *handler.AuthHandler
	SMTPService smtp.SMTPService
}

func NewAuthContainer(rdb *redis.Client, cfg *config.Config, db *gorm.DB, rabbitChan *amqp091.Channel) *AuthModule {
	mailer := smtp.NewSMTPService(cfg)
	authRepo := repoImpl.NewAuthRepository(rdb, cfg)
	userRepo := repoImpl.NewUserRepository(db)
	profileRepo := repoImpl.NewProfileRepository(db)
	authService := serviceImpl.NewAuthService(userRepo, authRepo, profileRepo, rabbitChan, cfg)
	userService := serviceImpl.NewUserService(userRepo, profileRepo)
	authHandler := handler.NewAuthHandler(authService, userService, cfg)

	return &AuthModule{
		authHandler,
		mailer,
	}
}
