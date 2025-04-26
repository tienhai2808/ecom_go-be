package common


import (
	"backend/internal/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppContext struct {
	DB          *gorm.DB
	Redis       *redis.Client
	Config      *config.AppConfig
}