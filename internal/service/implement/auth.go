package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/common"
	"github.com/tienhai2808/ecom_go/internal/dto"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/rabbitmq"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/response"
	"github.com/tienhai2808/ecom_go/internal/security"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/utils"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type authServiceImpl struct {
	userRepository    repository.UserRepository
	authRepository    repository.AuthRepository
	profileRepository repository.ProfileRepository
	rabbitChan        *amqp091.Channel
	config            *config.Config
}

func NewAuthService(userRepository repository.UserRepository, authRepository repository.AuthRepository, profileRepository repository.ProfileRepository, rabbitChan *amqp091.Channel, config *config.Config) service.AuthService {
	return &authServiceImpl{
		userRepository:    userRepository,
		authRepository:    authRepository,
		profileRepository: profileRepository,
		rabbitChan:        rabbitChan,
		config:            config,
	}
}

func (s *authServiceImpl) Signup(ctx context.Context, req request.SignupRequest) (string, error) {
	exists, err := s.userRepository.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return "", customErr.ErrEmailExists
	}

	exists, err = s.userRepository.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return "", customErr.ErrUsernameExists
	}

	otp := utils.GenerateOtp(5)
	registrationToken := uuid.NewString()

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	regData := dto.RegistrationData{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Otp:      otp,
		Attempts: 0,
	}

	if err = s.authRepository.AddRegistrationData(ctx, registrationToken, regData, 3*time.Minute); err != nil {
		return "", fmt.Errorf("lưu dữ liệu đăng ký thất bại: %w", err)
	}

	emailMsg := dto.EmailMessage{
		To:      req.Email,
		Subject: "Mã xác nhận Đăng ký tài khoản",
		Body:    fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp),
	}

	go func(msg dto.EmailMessage) {
		body, _ := json.Marshal(emailMsg)
		if err = rabbitmq.PublishMessage(s.rabbitChan, common.Exchange, common.RoutingKey, body); err != nil {
			log.Printf("publish email msg thất bại: %v", err)
		}
	}(emailMsg)

	return registrationToken, nil
}

func (s *authServiceImpl) VerifySignup(ctx context.Context, req request.VerifySignupRequest) (*response.AuthResponse, string, string, error) {
	regData, err := s.authRepository.GetRegistrationData(ctx, req.RegistrationToken)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy dữ liệu đăng ký thất bại: %w", err)
	}

	if regData == nil {
		return nil, "", "", customErr.ErrKeyNotFound
	}

	if regData.Attempts >= 3 {
		if err = s.authRepository.DeleteAuthData(ctx, "signup", req.RegistrationToken); err != nil {
			return nil, "", "", fmt.Errorf("xóa dữ liệu đăng ký thất bại: %w", err)
		}
		return nil, "", "", customErr.ErrTooManyAttempts
	}

	regData.Attempts++
	if err = s.authRepository.UpdateRegistrationData(ctx, req.RegistrationToken, *regData, 3*time.Minute); err != nil {
		return nil, "", "", fmt.Errorf("cập nhật dữ liệu đăng ký thất bại: %w", err)
	}

	if regData.Otp != req.Otp {
		return nil, "", "", customErr.ErrInvalidOTP
	}

	exists, err := s.userRepository.ExistsByEmail(ctx, regData.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return nil, "", "", customErr.ErrEmailExists
	}

	exists, err = s.userRepository.ExistsByUsername(ctx, regData.Username)
	if err != nil {
		return nil, "", "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return nil, "", "", customErr.ErrUsernameExists
	}

	newUser := &model.User{
		ID:       uuid.NewString(),
		Username: regData.Username,
		Email:    regData.Email,
		Password: regData.Password,
		Profile: &model.Profile{
			ID: uuid.NewString(),
		},
	}

	if err := s.userRepository.Create(ctx, newUser); err != nil {
		return nil, "", "", fmt.Errorf("tạo người dùng thất bại: %w", err)
	}

	accessToken, err := security.GenerateToken(newUser.ID, string(newUser.Role), 15*time.Minute, s.config.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo access_token thất bại: %w", err)
	}

	refreshToken, err := security.GenerateToken(newUser.ID, string(newUser.Role), 7*24*time.Hour, s.config.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo refresh_token thất bại: %w", err)
	}

	if err = s.authRepository.DeleteAuthData(ctx, "signup", req.RegistrationToken); err != nil {
		return nil, "", "", fmt.Errorf("xóa dữ liệu đăng ký thất bại: %w", err)
	}

	return s.ConvertToDto(newUser), accessToken, refreshToken, nil
}

func (s *authServiceImpl) Signin(ctx context.Context, req request.SigninRequest) (*response.AuthResponse, string, string, error) {
	user, err := s.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if user == nil {
		return nil, "", "", customErr.ErrUserNotFound
	}

	isCorrectPassword, err := utils.VerifyPassword(user.Password, req.Password)
	if err != nil || !isCorrectPassword {
		return nil, "", "", customErr.ErrIncorrectPassword
	}

	accessToken, err := security.GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.config.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo access_token thất bại: %w", err)
	}

	refreshToken, err := security.GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.config.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo refresh_token thất bại: %w", err)
	}

	return s.ConvertToDto(user), accessToken, refreshToken, nil
}

func (s *authServiceImpl) ForgotPassword(ctx context.Context, req request.ForgotPasswordRequest) (string, error) {
	exists, err := s.userRepository.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if !exists {
		return "", customErr.ErrUserNotFound
	}

	otp := utils.GenerateOtp(5)
	forgotPasswordToken := uuid.NewString()

	forgData := dto.ForgotPasswordData{
		Email:    req.Email,
		Otp:      otp,
		Attempts: 0,
	}

	if err = s.authRepository.AddForgotPasswordData(ctx, forgotPasswordToken, forgData, 3*time.Minute); err != nil {
		return "", fmt.Errorf("lưu dữ liệu quên mật khẩu thất bại: %w", err)
	}

	emailMsg := dto.EmailMessage{
		To:      req.Email,
		Subject: "Mã xác nhận Quên mật khẩu",
		Body:    fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp),
	}

	go func(msg dto.EmailMessage) {
		body, _ := json.Marshal(emailMsg)
		if err = rabbitmq.PublishMessage(s.rabbitChan, common.Exchange, common.RoutingKey, body); err != nil {
			log.Printf("publish email msg thất bại: %v", err)
		}
	}(emailMsg)

	return forgotPasswordToken, nil
}

func (s *authServiceImpl) VerifyForgotPassword(ctx context.Context, req request.VerifyForgotPasswordRequest) (string, error) {
	forgData, err := s.authRepository.GetForgotPasswordData(ctx, req.ForgotPasswordToken)
	if err != nil {
		return "", fmt.Errorf("lấy dữ liệu quên mật khẩu thất bại: %w", err)
	}

	if forgData == nil {
		return "", customErr.ErrKeyNotFound
	}

	if forgData.Attempts >= 3 {
		if err = s.authRepository.DeleteAuthData(ctx, "forgot-password", req.ForgotPasswordToken); err != nil {
			return "", fmt.Errorf("xóa dữ liệu quên mật khẩu thất bại: %w", err)
		}
		return "", customErr.ErrTooManyAttempts
	}

	if forgData.Otp != req.Otp {
		return "", customErr.ErrInvalidOTP
	}

	resetPasswordToken := uuid.NewString()

	if err = s.authRepository.AddResetPasswordData(ctx, resetPasswordToken, forgData.Email, 3*time.Minute); err != nil {
		return "", fmt.Errorf("lưu dữ liệu thay đổi mật khẩu thất bại: %w", err)
	}

	if err = s.authRepository.DeleteAuthData(ctx, "forgot-password", req.ForgotPasswordToken); err != nil {
		return "", fmt.Errorf("xóa dữ liệu quên mật khẩu thất bại: %w", err)
	}

	return resetPasswordToken, nil
}

func (s *authServiceImpl) ResetPassword(ctx context.Context, req request.ResetPasswordRequest) (*response.AuthResponse, string, string, error) {
	email, err := s.authRepository.GetResetPasswordData(ctx, req.ResetPasswordToken)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy dữ liệu làm mới mật khẩu thất bại: %w", err)
	}

	if email == "" {
		return nil, "", "", customErr.ErrKeyNotFound
	}

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if user == nil {
		return nil, "", "", customErr.ErrUserNotFound
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	if err = s.userRepository.UpdateUserPasswordByID(ctx, user.ID, hashedPassword); err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return nil, "", "", err
		}
		return nil, "", "", fmt.Errorf("cập nhật mật khẩu thất bại: %w", err)
	}

	accessToken, err := security.GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.config.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo access_token thất bại: %w", err)
	}

	refreshToken, err := security.GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.config.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo refresh_token thất bại: %w", err)
	}

	if err = s.authRepository.DeleteAuthData(ctx, "reset-password", req.ResetPasswordToken); err != nil {
		return nil, "", "", fmt.Errorf("xóa dữ liệu làm mới mật khẩu thất bại: %w", err)
	}

	return s.ConvertToDto(user), accessToken, refreshToken, nil
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, user *model.User, req request.ChangePasswordRequest) (*response.AuthResponse, string, string, error) {
	isCorrectPassword, err := utils.VerifyPassword(user.Password, req.OldPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("so sánh mật khẩu thất bại: %w", err)
	}

	if !isCorrectPassword {
		return nil, "", "", customErr.ErrIncorrectPassword
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("băm mật khẩu thất bại: %w", err)
	}

	if err = s.userRepository.UpdateUserPasswordByID(ctx, user.ID, hashedPassword); err != nil {
		if errors.Is(err, customErr.ErrUserNotFound) {
			return nil, "", "", err
		}
		return nil, "", "", fmt.Errorf("cập nhật mật khẩu thất bại: %w", err)
	}

	accessToken, err := security.GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.config.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo access_token thất bại: %w", err)
	}

	refreshToken, err := security.GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.config.App.JWTSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("tạo refresh_token thất bại: %w", err)
	}

	return s.ConvertToDto(user), accessToken, refreshToken, nil
}

func (s *authServiceImpl) UpdateUserProfile(ctx context.Context, user *model.User, req *request.UpdateProfileRequest) (*response.AuthResponse, error) {
	updateData := map[string]interface{}{}
	if req.FirstName != nil && *req.FirstName != user.Profile.FirstName {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != user.Profile.LastName {
		updateData["last_name"] = *req.LastName
	}
	if req.Gender != nil && *req.Gender != user.Profile.Gender {
		updateData["gender"] = *req.Gender
	}
	if req.DOB != nil && *req.DOB != user.Profile.DOB {
		updateData["dob"] = *req.DOB
	}
	if req.PhoneNumber != nil && *req.PhoneNumber != user.Profile.PhoneNumber {
		updateData["phone_number"] = *req.PhoneNumber
	}

	if len(updateData) > 0 {
		if err := s.profileRepository.UpdateProfileByUserID(ctx, user.ID, updateData); err != nil {
			if errors.Is(err, customErr.ErrUserProfileNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật thông tin người dùng thất bại: %w", err)
		}
	}

	updatedUser, err := s.userRepository.FindByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin người dùng thất bại: %w", err)
	}

	if updatedUser == nil {
		return nil, customErr.ErrUserNotFound
	}

	return s.ConvertToDto(updatedUser), nil
}

func (s *authServiceImpl) ConvertToDto(user *model.User) *response.AuthResponse {
	return &response.AuthResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Profile: response.ProfileResponse{
			FirstName:   user.Profile.FirstName,
			LastName:    user.Profile.LastName,
			PhoneNumber: user.Profile.PhoneNumber,
			Gender:      user.Profile.Gender,
			DOB:         user.Profile.DOB,
		},
	}
}
