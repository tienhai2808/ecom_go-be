package container

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/repository"
	repoImpl "backend/internal/repository/implement"
	svcImpl "backend/internal/service/implement"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthModule struct {
	UserRepository repository.UserRepository
	AuthHandler    handler.AuthHandler
}

func NewAuthContainer(redis *redis.Client, config *config.AppConfig, db *gorm.DB, rabbitChan *amqp091.Channel) *AuthModule {
	authRepo := repoImpl.NewAuthRepository(redis, config)
	userRepo := repoImpl.NewUserRepository(db)
	profileRepo := repoImpl.NewProfileRepository(db)
	addressRepo := repoImpl.NewAddressRepository(db)
	authService := svcImpl.NewAuthService(userRepo, authRepo, profileRepo, addressRepo, rabbitChan, config)
	authHandler := handler.NewAuthHandler(authService, config)
	return &AuthModule{
		UserRepository: userRepo,
		AuthHandler: *authHandler,
	}
}
