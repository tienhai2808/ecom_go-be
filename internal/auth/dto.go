package auth

import (
	"backend/internal/user"
	"time"
)

type TestRequest struct {
	Request string `json:"request"`
}

type SignupRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegistrationData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
	Attempts int    `json:"attempts"`
}

type VerifySignupRequest struct {
	RegistrationToken string `json:"registration_token" binding:"required,uuid4"`
	Otp               string `json:"otp" binding:"required,len=6,numeric"`
}

type SigninRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordData struct {
	Email    string `json:"email"`
	Otp      string `json:"otp"`
	Attempts int    `json:"attempts"`
}

type VerifyForgotPasswordRequest struct {
	ForgotPasswordToken string `json:"forgot_password_token" binding:"required,uuid4"`
	Otp                 string `json:"otp" binding:"required,len=6,numeric"`
}

type ResetPasswordRequest struct {
	ResetPasswordToken string `json:"reset_password_token" binding:"required,uuid4"`
	NewPassword        string `json:"new_password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UpdateInfoRequest struct {
	FirstName   *string          `json:"first_name,omitempty"`
	LastName    *string          `json:"last_name,omitempty"`
	Gender      *user.UserGender `json:"gender" binding:"omitempty,oneof=other male female"`
	DOB         *time.Time       `json:"dob,omitempty"`
	PhoneNumber *string          `json:"phone_number" binding:"omitempty,len=10"`
}

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type AddAddressRequest struct {
	FirstName   string `json:"first_name" binding:"omitempty"`
	LastName    string `json:"last_name" binding:"omitempty"`
	PhoneNumber string `json:"phone_number" binding:"required,len=10"`
	Address     string `json:"address" binding:"required"`
	Commune     string `json:"commune" binding:"required"`
	District    string `json:"district" binding:"required"`
	Province    string `json:"province" binding:"required"`
	IsDefault   bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	FirstName   *string `json:"first_name" binding:"omitempty"`
	LastName    *string `json:"last_name" binding:"omitempty"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,len=10"`
	Address     *string `json:"address" binding:"omitempty"`
	Commune     *string `json:"commune" binding:"omitempty,min=4"`
	District    *string `json:"district" binding:"omitempty,min=4"`
	Province    *string `json:"province" binding:"omitempty,min=4"`
	IsDefault   *bool   `json:"is_default" binding:"omitempty"`
}
