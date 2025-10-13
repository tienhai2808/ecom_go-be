package container

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	svcImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthModule struct {
	AuthHdl *handler.AuthHandler
}

func NewAuthContainer(rdb *redis.Client, cfg *config.Config, db *gorm.DB, rabbitChan *amqp091.Channel, sfg snowflake.SnowflakeGenerator) *AuthModule {
	authRepo := repoImpl.NewAuthRepository(rdb, cfg)
	userRepo := repoImpl.NewUserRepository(db)
	profileRepo := repoImpl.NewProfileRepository(db)
	authSvc := svcImpl.NewAuthService(userRepo, authRepo, profileRepo, rabbitChan, cfg, sfg)
	userSvc := svcImpl.NewUserService(userRepo, profileRepo, sfg)
	authHandler := handler.NewAuthHandler(authSvc, userSvc, cfg)

	return &AuthModule{authHandler}
}
