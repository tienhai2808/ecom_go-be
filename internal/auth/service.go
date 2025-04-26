package auth

import (
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/user"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Signup(req SignupRequest) (string, error)
	VerifySignup(req VerifySignupRequest) (*user.User, string, string, error)
	Signin(req SigninRequest) (*user.User, string, string, error)
	GetMe(userID string) (*user.User, error)
	ForgotPassword(req ForgotPasswordRequest) (string, error)
}

type service struct {
	repo Repository
	cfg  *config.AppConfig
}

func NewService(ctx *common.AppContext) Service {
	return &service{
		repo: NewRepository(ctx),
		cfg:  ctx.Config,
	}
}

func (s *service) Signup(req SignupRequest) (string, error) {
	exists, field, err := s.repo.CheckUserExists(req.Username, req.Email)
	if err != nil {
		return "", err
	}
	
	if exists {
		if field == "username" {
			return "", ErrUsernameExists
		}
		return "", ErrEmailExists
	}

	otp := GenerateOtp(5)
	registrationToken := uuid.NewString()

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	regData := RegistrationData{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Otp:      otp,
		Attempts: 0,
	}

	err = s.repo.StoreRegistrationData(registrationToken, regData, 3*time.Minute)
	if err != nil {
		return "", err
	}

	emailSender := common.NewSMTPSender(s.cfg)
	emailContent := fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp)
	err = emailSender.SendEmail(req.Email, "Mã xác nhận Đăng ký tài khoản", emailContent)
	if err != nil {
		s.repo.DeleteAuthData("signup", registrationToken)
		return "", fmt.Errorf("không thể gửi Email: %v", err)
	}

	return registrationToken, nil
}

func (s *service) VerifySignup(req VerifySignupRequest) (*user.User, string, string, error) {
	regData, err := s.repo.GetRegistrationData(req.RegistrationToken)
	if err != nil {
		return nil, "", "", ErrTokenExpired
	}

	if regData.Attempts >= 3 {
		s.repo.DeleteAuthData("signup", req.RegistrationToken)
		return nil, "", "", ErrTooManyAttempts
	}

	regData.Attempts++
	err = s.repo.UpdateRegistrationAttempts(req.RegistrationToken, *regData, 3*time.Minute)
	if err != nil {
		return nil, "", "", err
	}

	if regData.Otp != req.Otp {
		return nil, "", "", ErrInvalidOTP
	}

	exists, field, err := s.repo.CheckUserExists(regData.Username, regData.Email)
	if err != nil {
		return nil, "", "", err
	}
	
	if exists {
		if field == "username" {
			return nil, "", "", ErrUsernameExists
		}
		return nil, "", "", ErrEmailExists
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

	if err := s.repo.CreateUser(newUser); err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo người dùng: %v", err)
	}

	accessToken, err := GenerateToken(newUser.ID, string(newUser.Role), 15*time.Minute, s.cfg.App.JWTAccessSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo access_token: %v", err)
	}

	refreshToken, err := GenerateToken(newUser.ID, string(newUser.Role), 7*24*time.Hour, s.cfg.App.JWTRefreshSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo refresh_token: %v", err)
	}

	s.repo.DeleteAuthData("signup", req.RegistrationToken)

	return newUser, accessToken, refreshToken, nil
}

func (s *service) Signin(req SigninRequest) (*user.User, string, string, error) {
	// Get user by username
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, "", "", ErrUserNotFound
	}

	// Verify password
	isCorrectPassword, err := VerifyPassword(user.Password, req.Password)
	if err != nil || !isCorrectPassword {
		return nil, "", "", ErrIncorrectPassword
	}

	// Generate tokens
	accessToken, err := GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.cfg.App.JWTAccessSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo access_token: %v", err)
	}

	refreshToken, err := GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.cfg.App.JWTRefreshSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo refresh_token: %v", err)
	}

	return user, accessToken, refreshToken, nil
}

func (s *service) ForgotPassword(req ForgotPasswordRequest) (string, error) {
	exists, err := s.repo.CheckUserExistsByEmail(req.Email)
	if err != nil {
		return "", err
	}

	if exists {
		return "", ErrUserNotFound
	}

	otp := GenerateOtp(5)
	forgotPasswordToken := uuid.NewString()

	forgData := ForgotPasswordData{
		Email: req.Email,
		Otp: otp,
	}

	err = s.repo.StoreForgotPasswordData(forgotPasswordToken,  forgData, 3*time.Minute)
	if err != nil {
		return "", err
	}

	emailSender := common.NewSMTPSender(s.cfg)
	emailContent := fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp)
	err = emailSender.SendEmail(req.Email, "Mã xác nhận Quên mật khẩu", emailContent)
	if err != nil {
		s.repo.DeleteAuthData("forgot-password", forgotPasswordToken)
		return "", fmt.Errorf("không thể gửi Email: %v", err)
	}

	return forgotPasswordToken, nil
}

func (s *service) GetMe(userID string) (*user.User, error) {
	return s.repo.GetUserByID(userID)
}