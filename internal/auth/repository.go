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
	GetForgotPasswordData(token string) (*ForgotPasswordData, error)
	GetResetPasswordData(token string) (string, error)
	UpdateRegistrationAttempts(token string, data RegistrationData, ttl time.Duration) error
	UpdateForgotPasswordAttempts(token string, data ForgotPasswordData, ttl time.Duration) error
	DeleteAuthData(name, token string) error
	CreateUser(userData *user.User) error
	GetUserByUsername(username string) (*user.User, error)
	GetUserByID(id string) (*user.User, error)
	GetUserByEmail(email string) (*user.User, error)
	CheckUserExistsByEmail(email string) (bool, error)
	CheckUserExistsByUsername(username string) (bool, error)
	StoreForgotPasswordData(token string, data ForgotPasswordData, ttl time.Duration) error
	StoreResetPasswordData(token string, email string, ttl time.Duration) error
	UpdateUserPassword(userID string, hashedPassword string) error
	UpdateUserProfile(user *user.User, updateData map[string]interface{}) error
	UpdateUserInfo(user *user.User, updateData map[string]interface{}) error
	CreateAddress(addressData *user.Address) error
	UnsetDefaultAddress(userID string) error
	GetAddressByID(id string) (*user.Address, error)
	UpdateAddress(user *user.Address, updateData map[string]interface{}) error
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

func (r *repository) CheckUserExistsByUsername(username string) (bool, error) {
	var existingUser user.User

	err := r.db.Where("username = ?", username).First(&existingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
	if err := r.db.Preload("Profile").Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("không thể truy vấn người dùng: %v", err)
	}
	return &u, nil
}

func (r *repository) GetUserByID(id string) (*user.User, error) {
	var u user.User
	if err := r.db.Preload("Profile").Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("không thể truy vấn người dùng: %v", err)
	}
	return &u, nil
}

func (r *repository) GetUserByEmail(email string) (*user.User, error) {
	var u user.User
	if err := r.db.Preload("Profile").Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("không thể truy vấn người dùng: %v", err)
	}
	return &u, nil
}

func (r *repository) UpdateUserPassword(userID string, hashedPassword string) error {
	result := r.db.Model(&user.User{}).Where("id = ?", userID).Update("password", hashedPassword)
	if result.Error != nil {
		return fmt.Errorf("không thể cập nhật mật khẩu: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *repository) UpdateUserProfile(user *user.User, updateData map[string]interface{}) error {
	if err := r.db.Model(&user.Profile).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateUserInfo(user *user.User, updateData map[string]interface{}) error {
	if err := r.db.Model(&user).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) CheckUserExistsByEmail(email string) (bool, error) {
	var existingUser user.User

	err := r.db.Where("email = ?", email).First(&existingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repository) StoreForgotPasswordData(token string, data ForgotPasswordData, ttl time.Duration) error {
	forgDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu quên mật khẩu: %v", err)
	}

	ctxB := context.Background()
	redisKey := fmt.Sprintf("%s:forgot-password:%s", r.cfg.App.Name, token)
	err = r.redis.Set(ctxB, redisKey, forgDataJSON, ttl).Err()
	if err != nil {
		return fmt.Errorf("không thể lưu dữ liệu vào Redis: %v", err)
	}

	return nil
}

func (r *repository) StoreResetPasswordData(token string, email string, ttl time.Duration) error {
	ctxB := context.Background()
	redisKey := fmt.Sprintf("%s:reset-password:%s", r.cfg.App.Name, token)
	err := r.redis.Set(ctxB, redisKey, email, ttl).Err()
	if err != nil {
		return fmt.Errorf("không thể lưu dữ liệu vào Redis: %v", err)
	}

	return nil
}

func (r *repository) GetForgotPasswordData(token string) (*ForgotPasswordData, error) {
	redisKey := fmt.Sprintf("%s:forgot-password:%s", r.cfg.App.Name, token)
	ctxB := context.Background()

	forgDataJSON, err := r.redis.Get(ctxB, redisKey).Result()
	if err == redis.Nil {
		return nil, ErrTokenExpired
	}
	if err != nil {
		return nil, fmt.Errorf("không thể lấy dữ liệu từ Redis: %v", err)
	}

	var forgData ForgotPasswordData
	if err := json.Unmarshal([]byte(forgDataJSON), &forgData); err != nil {
		return nil, fmt.Errorf("không thể giải mã dữ liệu quên mật khẩu: %v", err)
	}

	return &forgData, nil
}

func (r *repository) UpdateForgotPasswordAttempts(token string, data ForgotPasswordData, ttl time.Duration) error {
	forgDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("không thể mã hóa dữ liệu quên mật khẩu: %v", err)
	}

	ctxB := context.Background()
	redisKey := fmt.Sprintf("%s:forgot-password:%s", r.cfg.App.Name, token)
	return r.redis.Set(ctxB, redisKey, forgDataJSON, ttl).Err()
}

func (r *repository) GetResetPasswordData(token string) (string, error) {
	redisKey := fmt.Sprintf("%s:reset-password:%s", r.cfg.App.Name, token)
	ctxB := context.Background()

	email, err := r.redis.Get(ctxB, redisKey).Result()
	if err == redis.Nil {
		return "", ErrTokenExpired
	}
	if err != nil {
		return "", fmt.Errorf("không thể lấy dữ liệu từ Redis: %v", err)
	}

	return email, nil
}

func (r *repository) CreateAddress(addressData *user.Address) error {
	return r.db.Create(addressData).Error
}

func (r *repository) UnsetDefaultAddress(userID string) error {
	return r.db.Model(&user.Address{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error
}

func (r *repository) GetAddressByID(id string) (*user.Address, error) {
	var a user.Address
	if err := r.db.Where("id = ?", id).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("không thể truy vấn địa chỉ: %v", err)
	}
	return &a, nil
}

func (r *repository) UpdateAddress(address *user.Address, updateData map[string]interface{}) error {
	if err := r.db.Model(&address).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}
