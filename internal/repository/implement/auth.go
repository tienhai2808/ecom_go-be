package implement

import (
	"backend/internal/config"
	"backend/internal/dto"
	"backend/internal/errors"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type authRepositoryImpl struct {
	redis  *redis.Client
	config *config.AppConfig
}

func NewAuthRepository(redis *redis.Client, config *config.AppConfig) repository.AuthRepository {
	return &authRepositoryImpl{
		redis:  redis,
		config: config,
	}
}

func (r *authRepositoryImpl) AddRegistrationData(token string, data dto.RegistrationData, ttl time.Duration) error {
	regData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu đăng ký: %w", err)
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("%s:signup:%s", r.config.App.Name, token)

	if err := r.redis.Set(ctx, redisKey, regData, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) DeleteAuthData(name, token string) error {
	redisKey := fmt.Sprintf("%s:%s:%s", r.config.App.Name, name, token)
	ctx := context.Background()

	if err := r.redis.Del(ctx, redisKey).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) GetRegistrationData(token string) (*dto.RegistrationData, error) {
	redisKey := fmt.Sprintf("%s:signup:%s", r.config.App.Name, token)
	ctx := context.Background()

	regDataJSON, err := r.redis.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, errors.ErrKeyNotFound
	} else if err != nil {
		return nil, fmt.Errorf("lỗi lấy dữ liệu từ redis: %w", err)
	}

	var regData dto.RegistrationData
	if err := json.Unmarshal([]byte(regDataJSON), &regData); err != nil {
		return nil, fmt.Errorf("không thể giải mã dữ liệu đăng ký: %w", err)
	}

	return &regData, nil
}

func (r *authRepositoryImpl) UpdateRegistrationData(token string, data dto.RegistrationData, ttl time.Duration) error {
	regDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu đăng ký: %w", err)
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("%s:signup:%s", r.config.App.Name, token)
	if err := r.redis.Set(ctx, redisKey, regDataJSON, ttl).Err(); err != nil {
		return err
	}

	return nil
}
