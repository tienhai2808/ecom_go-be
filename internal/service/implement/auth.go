package implement

import (
	"backend/internal/config"
	"backend/internal/dto"
	"backend/internal/errors"
	"backend/internal/model"
	"backend/internal/mq"
	"backend/internal/repository"
	"backend/internal/request"
	"backend/internal/service"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type authServiceImpl struct {
	userRepository repository.UserRepository
	authRepository repository.AuthRepository
	rabbitChan     *amqp091.Channel
	config         *config.AppConfig
}

func NewAuthService(userRepository repository.UserRepository, authRepository repository.AuthRepository, rabbitChan *amqp091.Channel, config *config.AppConfig) service.AuthService {
	return &authServiceImpl{
		userRepository: userRepository,
		authRepository: authRepository,
		rabbitChan:     rabbitChan,
		config:         config,
	}
}

func (s *authServiceImpl) Signup(req request.SignupRequest) (string, error) {
	exists, err := s.userRepository.CheckUserExistsByEmail(req.Email)
	if err != nil {
		return "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return "", errors.ErrEmailExists
	}

	exists, err = s.userRepository.CheckUserExistsByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return "", errors.ErrUsernameExists
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

	if err = s.authRepository.AddRegistrationData(registrationToken, regData, 3*time.Minute); err != nil {
		return "", fmt.Errorf("lưu dữ liệu đăng ký thất bại: %w", err)
	}

	emailMsg := dto.EmailMessage{
		To:      req.Email,
		Subject: "Mã xác nhận Đăng ký tài khoản",
		Body:    fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp),
	}

	body, err := json.Marshal(emailMsg)
	if err != nil {
		if err = s.authRepository.DeleteAuthData("signup", registrationToken); err != nil {
			return "", fmt.Errorf("xóa dữ liệu đăng ký thất bại: %w", err)
		}
		return "", fmt.Errorf("lỗi mã hóa email message: %w", err)
	}

	if err = mq.PublishMessage(s.rabbitChan, "", "email_queue", body); err != nil {
		if err = s.authRepository.DeleteAuthData("signup", registrationToken); err != nil {
			return "", fmt.Errorf("xóa dữ liệu đăng ký thất bại: %w", err)
		}
		return "", fmt.Errorf("không thể publish message: %w", err)
	}

	return registrationToken, nil
}

func (s *authServiceImpl) VerifySignup(req request.VerifySignupRequest) (*model.User, string, string, error) {
	regData, err := s.authRepository.GetRegistrationData(req.RegistrationToken)
	if err != nil {
		return nil, "", "", fmt.Errorf("lấy dữ liệu đăng ký thất bại: %w", err)
	}

	if regData.Attempts >= 3 {
		if err = s.authRepository.DeleteAuthData("signup", req.RegistrationToken); err != nil {
			return nil, "", "", fmt.Errorf("xóa dữ liệu đăng ký thất bại: %w", err)
		}
		return nil, "", "", errors.ErrTooManyAttempts
	}

	regData.Attempts++
	if err = s.authRepository.UpdateRegistrationData(req.RegistrationToken, *regData, 3*time.Minute); err != nil {
		return nil, "", "", fmt.Errorf("cập nhật dữ liệu đăng ký thất bại: %w", err)
	}

	if regData.Otp != req.Otp {
		return nil, "", "", errors.ErrInvalidOTP
	}

	exists, err := s.userRepository.CheckUserExistsByEmail(regData.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return nil, "", "", errors.ErrEmailExists
	}

	exists, err = s.userRepository.CheckUserExistsByUsername(regData.Username)
	if err != nil {
		return nil, "", "", fmt.Errorf("kiểm tra người dùng tồn tại thất bại: %w", err)
	}

	if exists {
		return nil, "", "", errors.ErrUsernameExists
	}

	newUser := &model.User{
		ID:       uuid.NewString(),
		Username: regData.Username,
		Email:    regData.Email,
		Password: regData.Password,
		Profile: model.Profile{
			ID: uuid.NewString(),
		},
	}

	if err := s.userRepository.CreateUser(newUser); err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo người dùng: %w", err)
	}

	accessToken, err := utils.GenerateToken(newUser.ID, string(newUser.Role), 15*time.Minute, s.config.App.JWTAccessSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo access_token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(newUser.ID, string(newUser.Role), 7*24*time.Hour, s.config.App.JWTRefreshSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo refresh_token: %w", err)
	}

	if err = s.authRepository.DeleteAuthData("signup", req.RegistrationToken); err != nil {
		return nil, "", "", fmt.Errorf("xóa dữ liệu đăng ký thất bại: %w", err)
	}

	return newUser, accessToken, refreshToken, nil
}

func (s *authServiceImpl) Signin(req request.SigninRequest) (*model.User, string, string, error) {
	user, err := s.userRepository.GetUserByUsername(req.Username)
	if err != nil {
		return nil, "", "", fmt.Errorf("đăng nhập thất bại: %w", err)
	}

	isCorrectPassword, err := utils.VerifyPassword(user.Password, req.Password)
	if err != nil || !isCorrectPassword {
		return nil, "", "", errors.ErrIncorrectPassword
	}

	accessToken, err := utils.GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.config.App.JWTAccessSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo access_token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.config.App.JWTRefreshSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo refresh_token: %w", err)
	}

	return user, accessToken, refreshToken, nil
}
