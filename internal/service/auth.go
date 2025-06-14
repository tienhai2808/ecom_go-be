package service

import (
	"backend/internal/model"
	"backend/internal/request"
	"backend/internal/response"
	"context"
)

type AuthService interface {
	Signup(ctx context.Context, req request.SignupRequest) (string, error)

	VerifySignup(ctx context.Context, req request.VerifySignupRequest) (*response.AuthResponse, string, string, error)

	Signin(ctx context.Context, req request.SigninRequest) (*response.AuthResponse, string, string, error)

	ForgotPassword(ctx context.Context, req request.ForgotPasswordRequest) (string, error)

	VerifyForgotPassword(ctx context.Context, req request.VerifyForgotPasswordRequest) (string, error)

	ResetPassword(ctx context.Context, req request.ResetPasswordRequest) (*response.AuthResponse, string, string, error)

	ChangePassword(ctx context.Context, user *model.User, req request.ChangePasswordRequest) (*response.AuthResponse, string, string, error)

	UpdateUserProfile(ctx context.Context, user *model.User, req *request.UpdateProfileRequest) (*response.AuthResponse, error)

	ConvertToDto(user *model.User) *response.AuthResponse
}
