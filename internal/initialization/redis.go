package initialization

import (
	"context"
	"fmt"
	"github.com/tienhai2808/ecom_go/internal/config"
	
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
		DB:   1,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("kết nối tới Redis thất bại: %w", err)
	}

	return rdb, nil
}
