package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/response"
)

type AuthService interface {
	SignUp(ctx context.Context, req request.SignUpRequest) (string, error)

	VerifySignUp(ctx context.Context, req request.VerifySignUpRequest) (*response.AuthResponse, string, string, error)

	SignIn(ctx context.Context, req request.SignInRequest) (*response.AuthResponse, string, string, error)

	ForgotPassword(ctx context.Context, req request.ForgotPasswordRequest) (string, error)

	VerifyForgotPassword(ctx context.Context, req request.VerifyForgotPasswordRequest) (string, error)

	ResetPassword(ctx context.Context, req request.ResetPasswordRequest) (*response.AuthResponse, string, string, error)

	ChangePassword(ctx context.Context, user *model.User, req request.ChangePasswordRequest) (*response.AuthResponse, string, string, error)

	UpdateProfile(ctx context.Context, user *model.User, req *request.UpdateProfileRequest) (*response.AuthResponse, error)

	ConvertToDto(user *model.User) *response.AuthResponse
}
