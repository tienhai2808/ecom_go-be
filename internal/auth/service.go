package auth

import (
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/user"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Service defines the auth service interface
type Service interface {
	Signup(req SignupRequest) (string, error)
	VerifySignup(req VerifySignupRequest) (*user.User, string, string, error)
	Signin(req SigninRequest) (*user.User, string, string, error)
	GetMe(userID string) (*user.User, error)
}

// service implements Service interface
type service struct {
	repo Repository
	cfg  *config.AppConfig
}

// NewService creates a new auth service
func NewService(ctx *common.AppContext) Service {
	return &service{
		repo: NewRepository(ctx),
		cfg:  ctx.Config,
	}
}

// Signup handles user registration
func (s *service) Signup(req SignupRequest) (string, error) {
	// Check if user exists
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

	// Generate OTP and registration token
	otp := GenerateOtp(5)
	registrationToken := uuid.NewString()
	
	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	// Prepare registration data
	regData := RegistrationData{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		OTP:      otp,
		Attempts: 0,
	}

	// Store registration data
	err = s.repo.StoreRegistrationData(registrationToken, regData, 3*time.Minute)
	if err != nil {
		return "", err
	}

	// Send email with OTP
	emailSender := common.NewSMTPSender(s.cfg)
	emailContent := fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp)
	err = emailSender.SendEmail(req.Email, "Mã xác nhận Đăng ký tài khoản", emailContent)
	if err != nil {
		s.repo.DeleteRegistrationData(registrationToken)
		return "", fmt.Errorf("không thể gửi Email: %v", err)
	}

	return registrationToken, nil
}

// VerifySignup verifies user registration
func (s *service) VerifySignup(req VerifySignupRequest) (*user.User, string, string, error) {
	// Get registration data
	regData, err := s.repo.GetRegistrationData(req.RegistrationToken)
	if err != nil {
		return nil, "", "", ErrTokenExpired
	}

	// Check attempts
	if regData.Attempts >= 3 {
		s.repo.DeleteRegistrationData(req.RegistrationToken)
		return nil, "", "", ErrTooManyAttempts
	}

	// Increment attempts
	regData.Attempts++
	err = s.repo.UpdateRegistrationAttempts(req.RegistrationToken, *regData, 3*time.Minute)
	if err != nil {
		return nil, "", "", err
	}

	// Verify OTP
	if regData.OTP != req.Otp {
		return nil, "", "", ErrInvalidOTP
	}

	// Check if user exists again (for race conditions)
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

	// Create new user
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

	// Generate tokens
	accessToken, err := GenerateToken(newUser.ID, string(newUser.Role), 15*time.Minute, s.cfg.App.JWTAccessSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo access_token: %v", err)
	}

	refreshToken, err := GenerateToken(newUser.ID, string(newUser.Role), 7*24*time.Hour, s.cfg.App.JWTRefreshSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo refresh_token: %v", err)
	}

	// Clean up registration data
	s.repo.DeleteRegistrationData(req.RegistrationToken)

	return newUser, accessToken, refreshToken, nil
}

// Signin handles user login
func (s *service) Signin(req SigninRequest) (*user.User, string, string, error) {
	// Get user by username
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, "", "", ErrUsernameNotFound
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

// GetMe retrieves user information
func (s *service) GetMe(userID string) (*user.User, error) {
	return s.repo.GetUserByID(userID)
}