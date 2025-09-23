package implement

import (
	"backend/config"
	"backend/internal/dto"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type authRepositoryImpl struct {
	redis  *redis.Client
	config *config.Config
}

func NewAuthRepository(redis *redis.Client, config *config.Config) repository.AuthRepository {
	return &authRepositoryImpl{
		redis:  redis,
		config: config,
	}
}

func (r *authRepositoryImpl) AddRegistrationData(ctx context.Context, token string, data dto.RegistrationData, ttl time.Duration) error {
	regData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu đăng ký: %w", err)
	}

	redisKey := fmt.Sprintf("%s:signup:%s", r.config.App.Name, token)

	if err := r.redis.Set(ctx, redisKey, regData, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) DeleteAuthData(ctx context.Context, name, token string) error {
	redisKey := fmt.Sprintf("%s:%s:%s", r.config.App.Name, name, token)

	if err := r.redis.Del(ctx, redisKey).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) GetRegistrationData(ctx context.Context, token string) (*dto.RegistrationData, error) {
	redisKey := fmt.Sprintf("%s:signup:%s", r.config.App.Name, token)

	regDataJSON, err := r.redis.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("lỗi lấy dữ liệu từ redis: %w", err)
	}

	var regData dto.RegistrationData
	if err = json.Unmarshal([]byte(regDataJSON), &regData); err != nil {
		return nil, fmt.Errorf("giải mã dữ liệu đăng ký thất bại: %w", err)
	}

	return &regData, nil
}

func (r *authRepositoryImpl) UpdateRegistrationData(ctx context.Context, token string, data dto.RegistrationData, ttl time.Duration) error {
	regDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu đăng ký: %w", err)
	}

	redisKey := fmt.Sprintf("%s:signup:%s", r.config.App.Name, token)
	if err := r.redis.Set(ctx, redisKey, regDataJSON, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) AddForgotPasswordData(ctx context.Context, token string, data dto.ForgotPasswordData, ttl time.Duration) error {
	forgDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu quên mật khẩu: %w", err)
	}

	redisKey := fmt.Sprintf("%s:forgot-password:%s", r.config.App.Name, token)

	if err = r.redis.Set(ctx, redisKey, forgDataJSON, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) GetForgotPasswordData(ctx context.Context, token string) (*dto.ForgotPasswordData, error) {
	redisKey := fmt.Sprintf("%s:forgot-password:%s", r.config.App.Name, token)

	forgDataJSON, err := r.redis.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("lỗi lấy dữ liệu từ redis: %w", err)
	}

	var forgData dto.ForgotPasswordData
	if err = json.Unmarshal([]byte(forgDataJSON), &forgData); err != nil {
		return nil, fmt.Errorf("giải mã dữ liệu quên mật khẩu thất bại: %w", err)
	}

	return &forgData, nil
}

func (r *authRepositoryImpl) AddResetPasswordData(ctx context.Context, token, email string, ttl time.Duration) error {
	redisKey := fmt.Sprintf("%s:reset-password:%s", r.config.App.Name, token)

	if err := r.redis.Set(ctx, redisKey, email, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) GetResetPasswordData(ctx context.Context, token string) (string, error) {
	redisKey := fmt.Sprintf("%s:reset-password:%s", r.config.App.Name, token)

	email, err := r.redis.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("lỗi lấy dữ liệu từ redis: %w", err)
	}

	return email, nil
}
