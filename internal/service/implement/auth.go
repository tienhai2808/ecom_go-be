package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	mrand "math/rand"
	"time"

	"github.com/tienhai2808/ecom_go/internal/common"
	"github.com/tienhai2808/ecom_go/internal/config"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/mapper"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/rabbitmq"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/response"
	"github.com/tienhai2808/ecom_go/internal/security"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"github.com/tienhai2808/ecom_go/internal/types"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type authServiceImpl struct {
	userRepo    repository.UserRepository
	authRepo    repository.AuthRepository
	profileRepo repository.ProfileRepository
	rabbitChan  *amqp091.Channel
	cfg         *config.Config
	sfg         snowflake.SnowflakeGenerator
}

func NewAuthService(userRepo repository.UserRepository, authRepo repository.AuthRepository, profileRepo repository.ProfileRepository, rabbitChan *amqp091.Channel, cfg *config.Config, sfg snowflake.SnowflakeGenerator) service.AuthService {
	return &authServiceImpl{
		userRepo,
		authRepo,
		profileRepo,
		rabbitChan,
		cfg,
		sfg,
	}
}

func (s *authServiceImpl) SignUp(ctx context.Context, req request.SignUpRequest) (string, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return "", customErr.ErrEmailExists
	}

	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return "", customErr.ErrUsernameExists
	}

	otp := generateOtp(5)
	registrationToken := uuid.NewString()

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return "", fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	regData := types.RegistrationData{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Otp:      otp,
		Attempts: 0,
	}

	if err = s.authRepo.AddRegistrationData(ctx, registrationToken, regData, 3*time.Minute); err != nil {
		return "", fmt.Errorf("lưu dữ liệu đăng ký thất bại: %w", err)
	}

	emailMsg := types.SendEmailMessage{
		To:      req.Email,
		Subject: "Mã xác nhận Đăng ký tài khoản",
		Body:    fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp),
	}

	go func(msg types.SendEmailMessage) {
		body, _ := json.Marshal(emailMsg)
		if err = rabbitmq.PublishMessage(s.rabbitChan, common.ExchangeEmail, common.RoutingKeyEmailSend, body); err != nil {
			log.Printf("publish email msg thất bại: %v", err)
		}
	}(emailMsg)

	return registrationToken, nil
}

func (s *authServiceImpl) VerifySignUp(ctx context.Context, req request.VerifySignUpRequest) (*response.UserResponse, string, string, error) {
	regData, err := s.authRepo.GetRegistrationData(ctx, req.RegistrationToken)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy dữ liệu đăng ký thất bại: %w", err)
	}

	if regData == nil {
		return nil, "", "", customErr.ErrKeyNotFound
	}

	if regData.Attempts >= 3 {
		if err = s.authRepo.DeleteAuthData(ctx, "signup", req.RegistrationToken); err != nil {
			return nil, "", "", fmt.Errorf("xóa dữ liệu đăng ký thất bại: %w", err)
		}
		return nil, "", "", customErr.ErrTooManyAttempts
	}

	regData.Attempts++
	if err = s.authRepo.UpdateRegistrationData(ctx, req.RegistrationToken, *regData, 3*time.Minute); err != nil {
		return nil, "", "", fmt.Errorf("cập nhật dữ liệu đăng ký thất bại: %w", err)
	}

	if regData.Otp != req.Otp {
		return nil, "", "", customErr.ErrInvalidOTP
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, regData.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return nil, "", "", customErr.ErrEmailExists
	}

	exists, err = s.userRepo.ExistsByUsername(ctx, regData.Username)
	if err != nil {
		return nil, "", "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return nil, "", "", customErr.ErrUsernameExists
	}

	userID, err := s.sfg.NextID()
	if err != nil {
		return nil, "", "", err
	}
	profileID, err := s.sfg.NextID()
	if err != nil {
		return nil, "", "", err
	}

	newUser := &model.User{
		ID:       userID,
		Username: regData.Username,
		Email:    regData.Email,
		Password: regData.Password,
		Profile: &model.Profile{
			ID: profileID,
		},
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, "", "", fmt.Errorf("tạo người dùng thất bại: %w", err)
	}

	accessToken, err := security.GenerateToken(newUser.ID, string(newUser.Role), 15*time.Minute, s.cfg.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo access_token thất bại: %w", err)
	}

	refreshToken, err := security.GenerateToken(newUser.ID, string(newUser.Role), 7*24*time.Hour, s.cfg.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo refresh_token thất bại: %w", err)
	}

	if err = s.authRepo.DeleteAuthData(ctx, "signup", req.RegistrationToken); err != nil {
		return nil, "", "", fmt.Errorf("xóa dữ liệu đăng ký thất bại: %w", err)
	}

	return mapper.ToUserResponse(newUser), accessToken, refreshToken, nil
}

func (s *authServiceImpl) SignIn(ctx context.Context, req request.SignInRequest) (*response.UserResponse, string, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if user == nil {
		return nil, "", "", customErr.ErrUserNotFound
	}

	isCorrectPassword, err := security.VerifyPassword(user.Password, req.Password)
	if err != nil || !isCorrectPassword {
		return nil, "", "", customErr.ErrIncorrectPassword
	}

	accessToken, err := security.GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.cfg.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo access_token thất bại: %w", err)
	}

	refreshToken, err := security.GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.cfg.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo refresh_token thất bại: %w", err)
	}

	return mapper.ToUserResponse(user), accessToken, refreshToken, nil
}

func (s *authServiceImpl) ForgotPassword(ctx context.Context, req request.ForgotPasswordRequest) (string, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if !exists {
		return "", customErr.ErrUserNotFound
	}

	otp := generateOtp(5)
	forgotPasswordToken := uuid.NewString()

	forgData := types.ForgotPasswordData{
		Email:    req.Email,
		Otp:      otp,
		Attempts: 0,
	}

	if err = s.authRepo.AddForgotPasswordData(ctx, forgotPasswordToken, forgData, 3*time.Minute); err != nil {
		return "", fmt.Errorf("lưu dữ liệu quên mật khẩu thất bại: %w", err)
	}

	emailMsg := types.SendEmailMessage{
		To:      req.Email,
		Subject: "Mã xác nhận Quên mật khẩu",
		Body:    fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp),
	}

	go func(msg types.SendEmailMessage) {
		body, _ := json.Marshal(emailMsg)
		if err = rabbitmq.PublishMessage(s.rabbitChan, common.ExchangeEmail, common.RoutingKeyEmailSend, body); err != nil {
			log.Printf("publish email msg thất bại: %v", err)
		}
	}(emailMsg)

	return forgotPasswordToken, nil
}

func (s *authServiceImpl) VerifyForgotPassword(ctx context.Context, req request.VerifyForgotPasswordRequest) (string, error) {
	forgData, err := s.authRepo.GetForgotPasswordData(ctx, req.ForgotPasswordToken)
	if err != nil {
		return "", fmt.Errorf("lấy dữ liệu quên mật khẩu thất bại: %w", err)
	}

	if forgData == nil {
		return "", customErr.ErrKeyNotFound
	}

	if forgData.Attempts >= 3 {
		if err = s.authRepo.DeleteAuthData(ctx, "forgot-password", req.ForgotPasswordToken); err != nil {
			return "", fmt.Errorf("xóa dữ liệu quên mật khẩu thất bại: %w", err)
		}
		return "", customErr.ErrTooManyAttempts
	}

	if forgData.Otp != req.Otp {
		return "", customErr.ErrInvalidOTP
	}

	resetPasswordToken := uuid.NewString()

	if err = s.authRepo.AddResetPasswordData(ctx, resetPasswordToken, forgData.Email, 3*time.Minute); err != nil {
		return "", fmt.Errorf("lưu dữ liệu thay đổi mật khẩu thất bại: %w", err)
	}

	if err = s.authRepo.DeleteAuthData(ctx, "forgot-password", req.ForgotPasswordToken); err != nil {
		return "", fmt.Errorf("xóa dữ liệu quên mật khẩu thất bại: %w", err)
	}

	return resetPasswordToken, nil
}

func (s *authServiceImpl) ResetPassword(ctx context.Context, req request.ResetPasswordRequest) (*response.UserResponse, string, string, error) {
	email, err := s.authRepo.GetResetPasswordData(ctx, req.ResetPasswordToken)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy dữ liệu làm mới mật khẩu thất bại: %w", err)
	}

	if email == "" {
		return nil, "", "", customErr.ErrKeyNotFound
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if user == nil {
		return nil, "", "", customErr.ErrUserNotFound
	}

	hashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	if err = s.userRepo.Update(ctx, user.ID, map[string]any{"password": hashedPassword}); err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return nil, "", "", err
		}
		return nil, "", "", fmt.Errorf("cập nhật mật khẩu thất bại: %w", err)
	}

	accessToken, err := security.GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.cfg.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo access_token thất bại: %w", err)
	}

	refreshToken, err := security.GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.cfg.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo refresh_token thất bại: %w", err)
	}

	if err = s.authRepo.DeleteAuthData(ctx, "reset-password", req.ResetPasswordToken); err != nil {
		return nil, "", "", fmt.Errorf("xóa dữ liệu làm mới mật khẩu thất bại: %w", err)
	}

	return mapper.ToUserResponse(user), accessToken, refreshToken, nil
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, userID int64, req request.ChangePasswordRequest) (*response.UserResponse, string, string, error) {
	user, err := s.userRepo.FindByIDWithProfile(ctx, userID)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}
	if user == nil {
		return nil, "", "", customErr.ErrUserNotFound
	}

	isCorrectPassword, err := security.VerifyPassword(user.Password, req.OldPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("so sánh mật khẩu thất bại: %w", err)
	}

	if !isCorrectPassword {
		return nil, "", "", customErr.ErrIncorrectPassword
	}

	hashedPassword, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	if err = s.userRepo.Update(ctx, user.ID, map[string]any{"password": hashedPassword}); err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return nil, "", "", err
		}
		return nil, "", "", fmt.Errorf("cập nhật mật khẩu thất bại: %w", err)
	}

	accessToken, err := security.GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.cfg.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo access_token thất bại: %w", err)
	}

	refreshToken, err := security.GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.cfg.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo refresh_token thất bại: %w", err)
	}

	return mapper.ToUserResponse(user), accessToken, refreshToken, nil
}

func generateOtp(length int) string {
	min := int(math.Pow10(length))
	max := 9 * min
	otp := min + mrand.Intn(max)
	return fmt.Sprintf("%d", otp)
}
