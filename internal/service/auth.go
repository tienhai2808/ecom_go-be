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

	GetMe(ctx context.Context, id string) (*model.User, error)

	ForgotPassword(ctx context.Context, req request.ForgotPasswordRequest) (string, error)

	VerifyForgotPassword(ctx context.Context, req request.VerifyForgotPasswordRequest) (string, error)

	ResetPassword(ctx context.Context, req request.ResetPasswordRequest) (*model.User, string, string, error)

	ChangePassword(ctx context.Context, user *model.User, req request.ChangePasswordRequest) (*model.User, string, string, error)

	UpdateUserProfile(ctx context.Context, user *model.User, req *request.UpdateProfileRequest) (*model.User, error)

	GetUserAddresses(ctx context.Context, userID string) ([]*model.Address, error)

	AddUserAddress(ctx context.Context, userID string, req request.AddAddressRequest) (*model.Address, error)
}
