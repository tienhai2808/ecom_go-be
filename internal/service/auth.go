package service

import (
	"backend/internal/model"
	"backend/internal/request"
	"context"
)

type AuthService interface {
	Signup(ctx context.Context, req request.SignupRequest) (string, error)

	VerifySignup(ctx context.Context, req request.VerifySignupRequest) (*model.User, string, string, error)

	Signin(ctx context.Context, req request.SigninRequest) (*model.User, string, string, error)

	GetMe(ctx context.Context, username string) (*model.User, error)

	ForgotPassword(ctx context.Context, req request.ForgotPasswordRequest) (string, error)

	VerifyForgotPassword(ctx context.Context, req request.VerifyForgotPasswordRequest) (string, error)

	ResetPassword(ctx context.Context, req request.ResetPasswordRequest) (*model.User, string, string, error)

	ChangePassword(ctx context.Context, id string, req request.ChangePasswordRequest) (*model.User, string, string, error)

	UpdateUserProfile(ctx context.Context, id string, req *request.UpdateProfileRequest) (*model.User, error)
}
