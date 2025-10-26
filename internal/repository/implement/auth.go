package implement

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/types"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type authRepositoryImpl struct {
	rdb  *redis.Client
	cfg *config.Config
}

func NewAuthRepository(rdb *redis.Client, cfg *config.Config) repository.AuthRepository {
	return &authRepositoryImpl{
		rdb,
		cfg,
	}
}

func (r *authRepositoryImpl) AddRegistrationData(ctx context.Context, token string, data types.RegistrationData, ttl time.Duration) error {
	regData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("mã hóa dữ liệu đăng ký thất bại: %w", err)
	}

	redisKey := fmt.Sprintf("%s:signup:%s", r.cfg.App.Name, token)

	if err := r.rdb.Set(ctx, redisKey, regData, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) DeleteAuthData(ctx context.Context, name, token string) error {
	redisKey := fmt.Sprintf("%s:%s:%s", r.cfg.App.Name, name, token)

	if err := r.rdb.Del(ctx, redisKey).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) GetRegistrationData(ctx context.Context, token string) (*types.RegistrationData, error) {
	redisKey := fmt.Sprintf("%s:signup:%s", r.cfg.App.Name, token)

	regDataJSON, err := r.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("lấy dữ liệu từ redis thất bại: %w", err)
	}

	var regData types.RegistrationData
	if err = json.Unmarshal([]byte(regDataJSON), &regData); err != nil {
		return nil, fmt.Errorf("giải mã dữ liệu đăng ký thất bại: %w", err)
	}

	return &regData, nil
}

func (r *authRepositoryImpl) UpdateRegistrationData(ctx context.Context, token string, data types.RegistrationData, ttl time.Duration) error {
	regDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("mã hóa dữ liệu đăng ký thất bại: %w", err)
	}

	redisKey := fmt.Sprintf("%s:signup:%s", r.cfg.App.Name, token)
	if err := r.rdb.Set(ctx, redisKey, regDataJSON, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) AddForgotPasswordData(ctx context.Context, token string, data types.ForgotPasswordData, ttl time.Duration) error {
	forgDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("mã hóa dữ liệu quên mật khẩu thất bại: %w", err)
	}

	redisKey := fmt.Sprintf("%s:forgot-password:%s", r.cfg.App.Name, token)

	if err = r.rdb.Set(ctx, redisKey, forgDataJSON, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) GetForgotPasswordData(ctx context.Context, token string) (*types.ForgotPasswordData, error) {
	redisKey := fmt.Sprintf("%s:forgot-password:%s", r.cfg.App.Name, token)

	forgDataJSON, err := r.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("lấy dữ liệu từ redis thất bại: %w", err)
	}

	var forgData types.ForgotPasswordData
	if err = json.Unmarshal([]byte(forgDataJSON), &forgData); err != nil {
		return nil, fmt.Errorf("giải mã dữ liệu quên mật khẩu thất bại: %w", err)
	}

	return &forgData, nil
}

func (r *authRepositoryImpl) AddResetPasswordData(ctx context.Context, token, email string, ttl time.Duration) error {
	redisKey := fmt.Sprintf("%s:reset-password:%s", r.cfg.App.Name, token)

	if err := r.rdb.Set(ctx, redisKey, email, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRepositoryImpl) GetResetPasswordData(ctx context.Context, token string) (string, error) {
	redisKey := fmt.Sprintf("%s:reset-password:%s", r.cfg.App.Name, token)

	email, err := r.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("lấy dữ liệu từ redis thất bại: %w", err)
	}

	return email, nil
}
