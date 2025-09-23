package repository

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/dto"
	"time"
)

type AuthRepository interface {
	AddRegistrationData(ctx context.Context, token string, data dto.RegistrationData, ttl time.Duration) error

	DeleteAuthData(ctx context.Context, name, token string) error

	GetRegistrationData(ctx context.Context, token string) (*dto.RegistrationData, error)

	UpdateRegistrationData(ctx context.Context, token string, data dto.RegistrationData, ttl time.Duration) error

	AddForgotPasswordData(ctx context.Context, token string, data dto.ForgotPasswordData, ttl time.Duration) error

	GetForgotPasswordData(ctx context.Context, token string) (*dto.ForgotPasswordData, error)

	AddResetPasswordData(ctx context.Context, token, email string, ttl time.Duration) error

	GetResetPasswordData(ctx context.Context, token string) (string, error)
}
