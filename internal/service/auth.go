package service

import (
	"backend/internal/model"
	"backend/internal/request"
)

type AuthService interface {
	Signup(req request.SignupRequest) (string, error)
	VerifySignup(req request.VerifySignupRequest) (*model.User, string, string, error)
	Signin(req request.SigninRequest) (*model.User, string, string, error)
	GetMe(username string) (*model.User, error)
	ForgotPassword(req request.ForgotPasswordRequest) (string, error)
}
