package auth

import (
	"context"
	"e-commerce/internal/common"
	"e-commerce/internal/user"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RegistrationData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
	Attempts int `json:"attempts"`
}

func SignupService(ctx *common.AppContext, req SignupSchema) (string, error) {
	var existingUser user.User
	
	err := ctx.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error
	if err == nil {
		if existingUser.Username == req.Username {
			return "", ErrUsernameExists
		}
		if existingUser.Email == req.Email {
			return "", ErrEmailExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	otp := GenerateOtp(5)
	registrationToken:= uuid.NewString()
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	regData := RegistrationData{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		OTP:      otp,
		Attempts: 0,
	}

	regDataJSON, err := json.Marshal(regData)
	if err != nil {
		return "", fmt.Errorf("không thể mã hóa dữ liệu đăng ký: %v", err)
	}

	ctxB := context.Background()
	redisKey := fmt.Sprintf("%s:signup:%s", ctx.Config.App.Name, registrationToken)
	err = ctx.Redis.Set(ctxB, redisKey, regDataJSON, 3*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("không thể lưu dữ liệu vào Redis: %v", err) 
	}

	emailSender := common.NewSMTPSender(ctx.Config) 
	emailContent := fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp)
	err = emailSender.SendEmail(req.Email, "Mã xác nhận Đăng ký tài khoản", emailContent)
  if err != nil {
		ctx.Redis.Del(ctxB, redisKey)
		return "", fmt.Errorf("không thể gửi Email: %v", err)
	}

	return registrationToken, nil
}

func VerifySignupService(ctx *common.AppContext, req VerifySignupSchema) (*user.User, string, string, error) {
	redisKey := fmt.Sprintf("%s:signup:%s", ctx.Config.App.Name, req.RegistrationToken)
	ctxB := context.Background()

	regDataJSON, err := ctx.Redis.Get(ctxB, redisKey).Result()
	if err == redis.Nil {
		return nil, "", "", ErrTokenExpired
	}
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể lấy dữ liệu từ Redis: %v", err)
	}

	var regData RegistrationData
	if err := json.Unmarshal([]byte(regDataJSON), &regData); err != nil {
		return nil, "", "", fmt.Errorf("không thể giải mã dữ liệu đăng ký: %v", err)
	}


	if regData.Attempts >= 3 {
		ctx.Redis.Del(ctxB, redisKey)
		return nil, "", "", ErrTooManyAttempts
	}

	regData.Attempts++
	regDataJSONBytes, err := json.Marshal(regData)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể mã hóa dữ liệu đăng ký: %v", err)
	}
	regDataJSON = string(regDataJSONBytes)
	ctx.Redis.Set(ctxB, redisKey, regDataJSON, 3*time.Minute)

	if regData.OTP != req.Otp {
		return nil, "", "", ErrInvalidOTP
	}

	var existingUser user.User
	err = ctx.DB.Where("username = ? OR email = ?", regData.Username, regData.Email).First(&existingUser).Error
	if err == nil {
		if existingUser.Username == regData.Username {
			return nil, "", "", ErrUsernameExists
		}
		if existingUser.Email == regData.Email {
			return nil, "", "", ErrEmailExists
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", "", err
	}

	newUser := &user.User{
		ID:       uuid.NewString(),
		Username: regData.Username,
		Email:    regData.Email,
		Password: regData.Password,
		Profile: user.Profile{
			ID: uuid.NewString(),
		},
	}

	if err := ctx.DB.Create(newUser).Error; err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo người dùng: %v", err)
	}

	accessToken, err := GenerateToken(newUser.ID, string(newUser.Role), 15*time.Minute, ctx.Config.App.JWTAccessSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo access_token: %v", err)
	}

	refreshToken, err := GenerateToken(newUser.ID, string(newUser.Role), 7*24*time.Hour, ctx.Config.App.JWTRefreshSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo refresh_token: %v", err)
	}

	ctx.Redis.Del(ctxB, redisKey)

	return newUser, accessToken, refreshToken, nil
}

func GetMeService(ctx *common.AppContext, userID string) (*user.User, error) {
	var u user.User

	err := ctx.DB.Where("id = ?", userID).First(&u).Error; 
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("người dùng không tồn tại")
		}
		return nil, fmt.Errorf("không thể truy vấn người dùng: %v", err)
	}
	return &u, nil
}
