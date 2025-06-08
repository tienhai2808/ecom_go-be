package container

import (
	"backend/internal/config"
	"backend/internal/handler"
	repoImpl "backend/internal/repository/implement"
	serviceImpl "backend/internal/service/implement"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthModule struct {
	AuthHandler    handler.AuthHandler
}

func NewAuthContainer(redis *redis.Client, config *config.AppConfig, db *gorm.DB, rabbitChan *amqp091.Channel) *AuthModule {
	authRepo := repoImpl.NewAuthRepository(redis, config)
	userRepo := repoImpl.NewUserRepository(db)
	profileRepo := repoImpl.NewProfileRepository(db)
	authService := serviceImpl.NewAuthService(userRepo, authRepo, profileRepo, rabbitChan, config)
	userService := serviceImpl.NewUserService(userRepo, profileRepo)
	authHandler := handler.NewAuthHandler(authService, userService, config)
	
	return &AuthModule{
		AuthHandler: *authHandler,
	}
}
