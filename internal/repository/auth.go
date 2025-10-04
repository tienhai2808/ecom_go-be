package repository

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/types"
	"time"
)

type AuthRepository interface {
	AddRegistrationData(ctx context.Context, token string, data types.RegistrationData, ttl time.Duration) error

	DeleteAuthData(ctx context.Context, name, token string) error

	GetRegistrationData(ctx context.Context, token string) (*types.RegistrationData, error)

	UpdateRegistrationData(ctx context.Context, token string, data types.RegistrationData, ttl time.Duration) error

	AddForgotPasswordData(ctx context.Context, token string, data types.ForgotPasswordData, ttl time.Duration) error

	GetForgotPasswordData(ctx context.Context, token string) (*types.ForgotPasswordData, error)

	AddResetPasswordData(ctx context.Context, token, email string, ttl time.Duration) error

	GetResetPasswordData(ctx context.Context, token string) (string, error)
}
