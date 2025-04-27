package admin

import (
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/user"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	GetUsers() ([]user.User, error)
}

type repository struct {
	db    *gorm.DB
	redis *redis.Client
	cfg   *config.AppConfig
}

func NewRepository(ctx *common.AppContext) Repository {
	return &repository{
		db:    ctx.DB,
		redis: ctx.Redis,
		cfg:   ctx.Config,
	}
}

func (r *repository) GetUsers() ([]user.User, error) {
	var users []user.User

	err := r.db.Preload("Profile").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
