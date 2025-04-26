package auth

import (
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	CheckUserExists(username, email string) (bool, string, error)
	StoreRegistrationData(token string, data RegistrationData, ttl time.Duration) error
	GetRegistrationData(token string) (*RegistrationData, error)
	UpdateRegistrationAttempts(token string, data RegistrationData, ttl time.Duration) error
	DeleteAuthData(name, token string) error
	CreateUser(userData *user.User) error
	GetUserByUsername(username string) (*user.User, error)
	GetUserByID(id string) (*user.User, error)
	CheckUserExistsByEmail(email string) (bool, error)
	StoreForgotPasswordData(token string, data ForgotPasswordData, ttl time.Duration) error
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

func (r *repository) CheckUserExists(username, email string) (bool, string, error) {
	var existingUser user.User
	
	err := r.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error
	if err == nil {
		if existingUser.Username == username {
			return true, "username", nil
		}
		if existingUser.Email == email {
			return true, "email", nil
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, "", err
	}
	
	return false, "", nil
}

func (r *repository) StoreRegistrationData(token string, data RegistrationData, ttl time.Duration) error {
	regDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu đăng ký: %v", err)
	}

	ctxB := context.Background()
	redisKey := fmt.Sprintf("%s:signup:%s", r.cfg.App.Name, token)
	err = r.redis.Set(ctxB, redisKey, regDataJSON, ttl).Err()
	if err != nil {
		return fmt.Errorf("không thể lưu dữ liệu vào Redis: %v", err) 
	}
	
	return nil
}

func (r *repository) GetRegistrationData(token string) (*RegistrationData, error) {
	redisKey := fmt.Sprintf("%s:signup:%s", r.cfg.App.Name, token)
	ctxB := context.Background()

	regDataJSON, err := r.redis.Get(ctxB, redisKey).Result()
	if err == redis.Nil {
		return nil, ErrTokenExpired
	}
	if err != nil {
		return nil, fmt.Errorf("không thể lấy dữ liệu từ Redis: %v", err)
	}

	var regData RegistrationData
	if err := json.Unmarshal([]byte(regDataJSON), &regData); err != nil {
		return nil, fmt.Errorf("không thể giải mã dữ liệu đăng ký: %v", err)
	}
	
	return &regData, nil
}

func (r *repository) UpdateRegistrationAttempts(token string, data RegistrationData, ttl time.Duration) error {
	regDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu đăng ký: %v", err)
	}

	ctxB := context.Background()
	redisKey := fmt.Sprintf("%s:signup:%s", r.cfg.App.Name, token)
	return r.redis.Set(ctxB, redisKey, regDataJSON, ttl).Err()
}

func (r *repository) DeleteAuthData(name, token string) error {
	redisKey := fmt.Sprintf("%s:%s:%s", r.cfg.App.Name, name, token)
	ctxB := context.Background()
	return r.redis.Del(ctxB, redisKey).Err()
}

func (r *repository) CreateUser(userData *user.User) error {
	return r.db.Create(userData).Error
}

func (r *repository) GetUserByUsername(username string) (*user.User, error) {
	var u user.User
	err := r.db.Preload("Profile").Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("không thể truy vấn người dùng: %v", err)
	}
	return &u, nil
}

func (r *repository) GetUserByID(id string) (*user.User, error) {
	var u user.User
	err := r.db.Preload("Profile").Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("không thể truy vấn người dùng: %v", err)
	}
	return &u, nil
}

func (r *repository) CheckUserExistsByEmail(email string) (bool, error) {
	var existingUser user.User
	
	err := r.db.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		return false, err
	}
	return true, nil
}

func (r *repository) StoreForgotPasswordData(token string, data ForgotPasswordData, ttl time.Duration) error {
	regDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu quên mật khẩu: %v", err)
	}

	ctxB := context.Background()
	redisKey := fmt.Sprintf("%s:forgot-password:%s", r.cfg.App.Name, token)
	err = r.redis.Set(ctxB, redisKey, regDataJSON, ttl).Err()
	if err != nil {
		return fmt.Errorf("không thể lưu dữ liệu vào Redis: %v", err) 
	}
	
	return nil
}