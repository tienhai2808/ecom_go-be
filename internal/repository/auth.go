package repository

import (
	"backend/internal/dto"
	"time"
)

type AuthRepository interface {
	AddRegistrationData(token string, data dto.RegistrationData, ttl time.Duration) error
	DeleteAuthData(name, token string) error
	GetRegistrationData(token string) (*dto.RegistrationData, error)
	UpdateRegistrationData(token string, data dto.RegistrationData, ttl time.Duration) error
	AddForgotPasswordData(token string, data dto.ForgotPasswordData, ttl time.Duration) error
}
