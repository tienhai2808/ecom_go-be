package auth

import (
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/mq"
	"backend/internal/user"
	"encoding/json"
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
	VerifyForgotPassword(req VerifyForgotPasswordRequest) (string, error)
	ResetPassword(req ResetPasswordRequest) (*user.User, string, string, error)
	ChangePassword(userID string, req ChangePasswordRequest) (*user.User, string, string, error)
	UpdateInfo(userID string, req *UpdateInfoRequest) (*user.User, error)
	AddAddress(userID string, req AddAddressRequest) (*user.Address, error)
	UpdateAddress(userID, addressID string, req UpdateAddressRequest) (*user.Address, error)
}

type service struct {
	repo Repository
	cfg  *config.AppConfig
	ctx  *common.AppContext
}

func NewService(ctx *common.AppContext) Service {
	return &service{
		repo: NewRepository(ctx),
		cfg:  ctx.Config,
		ctx:  ctx,
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

	emailMsg := EmailMessage{
		To:      req.Email,
		Subject: "Mã xác nhận Đăng ký tài khoản",
		Body:    fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp),
	}

	body, err := json.Marshal(emailMsg)
	if err != nil {
		s.repo.DeleteAuthData("signup", registrationToken)
		return "", fmt.Errorf("lỗi khi marshal email message: %v", err)
	}

	err = mq.PublishMessage(s.ctx.RabbitChan, "", "email_queue", body)
	if err != nil {
		s.repo.DeleteAuthData("signup", registrationToken)
		return "", fmt.Errorf("không thể publish message: %v", err)
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
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, "", "", ErrUserNotFound
	}

	isCorrectPassword, err := VerifyPassword(user.Password, req.Password)
	if err != nil || !isCorrectPassword {
		return nil, "", "", ErrIncorrectPassword
	}

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

	if !exists {
		return "", ErrUserNotFound
	}

	otp := GenerateOtp(5)
	forgotPasswordToken := uuid.NewString()

	forgData := ForgotPasswordData{
		Email:    req.Email,
		Otp:      otp,
		Attempts: 0,
	}

	err = s.repo.StoreForgotPasswordData(forgotPasswordToken, forgData, 3*time.Minute)
	if err != nil {
		return "", err
	}

	emailMsg := EmailMessage{
		To:      req.Email,
		Subject: "Mã xác nhận Quên mật khẩu",
		Body:    fmt.Sprintf(`Đây là mã OTP của bạn, nó sẽ hết hạn sau 3 phút: <p style="text-align: center"><strong style="font-size: 18px; color: #333;">%s</strong></p>`, otp),
	}

	body, err := json.Marshal(emailMsg)
	if err != nil {
		s.repo.DeleteAuthData("forgot-password", forgotPasswordToken)
		return "", fmt.Errorf("lỗi khi marshal email message: %v", err)
	}

	err = mq.PublishMessage(s.ctx.RabbitChan, "", "email_queue", body)
	if err != nil {
		s.repo.DeleteAuthData("forgot-password", forgotPasswordToken)
		return "", fmt.Errorf("không thể publish message: %v", err)
	}

	return forgotPasswordToken, nil
}

func (s *service) VerifyForgotPassword(req VerifyForgotPasswordRequest) (string, error) {
	forgData, err := s.repo.GetForgotPasswordData(req.ForgotPasswordToken)
	if err != nil {
		return "", ErrTokenExpired
	}

	if forgData.Attempts >= 3 {
		s.repo.DeleteAuthData("forgot-password", req.ForgotPasswordToken)
		return "", ErrTooManyAttempts
	}

	forgData.Attempts++
	err = s.repo.UpdateForgotPasswordAttempts(req.ForgotPasswordToken, *forgData, 3*time.Minute)
	if err != nil {
		return "", err
	}

	if forgData.Otp != req.Otp {
		return "", ErrInvalidOTP
	}

	resetPasswordToken := uuid.NewString()

	err = s.repo.StoreResetPasswordData(resetPasswordToken, forgData.Email, 3*time.Minute)
	if err != nil {
		return "", err
	}

	s.repo.DeleteAuthData("forgot-password", req.ForgotPasswordToken)

	return resetPasswordToken, nil
}

func (s *service) ResetPassword(req ResetPasswordRequest) (*user.User, string, string, error) {
	email, err := s.repo.GetResetPasswordData(req.ResetPasswordToken)
	if err != nil {
		return nil, "", "", ErrTokenExpired
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, "", "", ErrUserNotFound
	}

	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể mã hóa mật khẩu: %v", err)
	}

	if err := s.repo.UpdateUserPassword(user.ID, hashedPassword); err != nil {
		return nil, "", "", fmt.Errorf("không thể cập nhật mật khẩu: %v", err)
	}

	accessToken, err := GenerateToken(user.ID, string(user.Role), 15*time.Minute, s.cfg.App.JWTAccessSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo access_token: %v", err)
	}

	refreshToken, err := GenerateToken(user.ID, string(user.Role), 7*24*time.Hour, s.cfg.App.JWTRefreshSecret)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể tạo refresh_token: %v", err)
	}

	s.repo.DeleteAuthData("reset-password", req.ResetPasswordToken)

	return user, accessToken, refreshToken, nil
}

func (s *service) GetMe(userID string) (*user.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *service) ChangePassword(userID string, req ChangePasswordRequest) (*user.User, string, string, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, "", "", ErrUserNotFound
	}

	isCorrectPassword, err := VerifyPassword(user.Password, req.OldPassword)
	if err != nil || !isCorrectPassword {
		return nil, "", "", ErrIncorrectPassword
	}

	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		return nil, "", "", fmt.Errorf("không thể mã hóa mật khẩu: %v", err)
	}

	if err := s.repo.UpdateUserPassword(userID, hashedPassword); err != nil {
		return nil, "", "", fmt.Errorf("không thể cập nhật mật khẩu: %v", err)
	}

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

func (s *service) UpdateInfo(userID string, req *UpdateInfoRequest) (*user.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	updateData := map[string]interface{}{}
	if req.FirstName != nil {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updateData["last_name"] = *req.LastName
	}
	if req.Gender != nil {
		updateData["gender"] = *req.Gender
	}
	if req.DOB != nil {
		updateData["dob"] = *req.DOB
	}
	if req.PhoneNumber != nil {
		updateData["phone_number"] = *req.PhoneNumber
	}

	if len(updateData) > 0 {
		if err := s.repo.UpdateUserProfile(user, updateData); err != nil {
			return nil, ErrUpdateFailed
		}
	}

	updatedUser, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return updatedUser, nil
}

func (s *service) AddAddress(userID string, req AddAddressRequest) (*user.Address, error) {
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	exists, err := s.repo.CheckDefaultAddressExists(userID)
	if err != nil {
		return nil, fmt.Errorf("lỗi kiểm tra địa chỉ mặc định: %v", err)
	}
	if exists && req.IsDefault {
		if err := s.repo.UnsetDefaultAddress(userID); err != nil {
			return nil, fmt.Errorf("lỗi cập nhật địa chỉ mặc định cũ: %v", err)
		}
	}
	if !exists && !req.IsDefault {
		req.IsDefault = true
	}

	newAddress := &user.Address{
		ID:          uuid.NewString(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Commune:     req.Commune,
		District:    req.District,
		Province:    req.Province,
		Address:     req.Address,
		IsDefault:   req.IsDefault,
		UserID:      userID,
	}
	if err := s.repo.CreateAddress(newAddress); err != nil {
		return nil, fmt.Errorf("không thể tạo địa chỉ: %v", err)
	}
	return newAddress, nil
}

func (s *service) UpdateAddress(userID, addressID string, req UpdateAddressRequest) (*user.Address, error) {
	address, err := s.repo.GetAddressByID(addressID)
	if err != nil {
		return nil, ErrAddressNotFound
	}

	if address.UserID != userID {
		return nil, ErrUnAuth
	}

	updateData := map[string]interface{}{}
	if req.FirstName != nil {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updateData["last_name"] = *req.LastName
	}
	if req.PhoneNumber != nil {
		updateData["phone_number"] = *req.PhoneNumber
	}
	if req.Address != nil {
		updateData["address"] = *req.Address
	}
	if req.Commune != nil {
		updateData["commune"] = *req.Commune
	}
	if req.District != nil {
		updateData["district"] = *req.District
	}
	if req.Province != nil {
		updateData["province"] = *req.Province
	}
	if req.IsDefault != nil && *req.IsDefault != address.IsDefault {
		if *req.IsDefault {
			if err := s.repo.UnsetDefaultAddress(userID); err != nil {
				return nil, fmt.Errorf("lỗi cập nhật địa chỉ mặc định cũ: %v", err)
			}
		} else {
			countAddress, err := s.repo.CountAddress(userID)
			if err != nil {
				return nil, fmt.Errorf("lỗi đếm địa chỉ: %v", err)
			}
			if countAddress > 1 {
				if err := s.repo.SetLatestDefaultAddress(userID, addressID); err != nil {
					return nil, fmt.Errorf("lỗi xét lại địa chỉ mặc định: %v", err)
				}
			} else {
				*req.IsDefault = true
			}
			updateData["is_default"] = *req.IsDefault
		}
	}

	if len(updateData) > 0 {
		if err := s.repo.UpdateAddress(address, updateData); err != nil {
			return nil, ErrUpdateFailed
		}
	}

	updatedAddress, err := s.repo.GetAddressByID(address.ID)
	if err != nil {
		return nil, ErrAddressNotFound
	}

	return updatedAddress, nil
}
