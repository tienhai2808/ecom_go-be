package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/response"
)

type AuthService interface {
	SignUp(ctx context.Context, req request.SignUpRequest) (string, error)

	VerifySignUp(ctx context.Context, req request.VerifySignUpRequest) (*response.UserResponse, string, string, error)

	SignIn(ctx context.Context, req request.SignInRequest) (*response.UserResponse, string, string, error)

	ForgotPassword(ctx context.Context, req request.ForgotPasswordRequest) (string, error)

	VerifyForgotPassword(ctx context.Context, req request.VerifyForgotPasswordRequest) (string, error)

	ResetPassword(ctx context.Context, req request.ResetPasswordRequest) (*response.UserResponse, string, string, error)

	ChangePassword(ctx context.Context, userID int64, req request.ChangePasswordRequest) (*response.UserResponse, string, string, error)
}
