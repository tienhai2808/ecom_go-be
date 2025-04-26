package cache

import (
	"context"
	"backend/internal/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func ConnectToRedis() (*redis.Client, error) {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("❤️ lỗi load config redis: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})

	if err := rdb.Ping(Ctx).Err(); err != nil {
		return nil, fmt.Errorf("❤️ không thể kết nối Redis: %v", err)
	}

	Rdb = rdb
	return rdb, nil
}



