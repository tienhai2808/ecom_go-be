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
	DeleteUSers(userIDs []string) (int64, error)
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

	if err := r.db.Preload("Profile").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) DeleteUSers(userIDs []string) (int64, error) {
	result := r.db.Where("id IN ?", userIDs).Delete(&user.User{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}