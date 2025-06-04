package container

import (
	"backend/internal/config"
	"backend/internal/handler"
	repository "backend/internal/repository/implement"
	service "backend/internal/service/implement"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthModule struct {
	AuthHandler handler.AuthHandler
}

func NewAuthContainer(redis *redis.Client, config *config.AppConfig, db *gorm.DB, rabbitChan *amqp091.Channel) *AuthModule {
	authRepo := repository.NewAuthRepository(redis, config)
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	authService := service.NewAuthService(userRepo, authRepo, profileRepo, rabbitChan, config)
	authHandler := handler.NewAuthHandler(authService, config)
	return &AuthModule{
		AuthHandler: *authHandler,
	}
}
