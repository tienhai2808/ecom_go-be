package admin

import (
	"backend/internal/auth"
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/user"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	GetAllUsers() ([]user.User, error)
	CreateUser(req CreateUserRequest) error
}

type service struct {
	repo     Repository
	authRepo auth.Repository
	cfg      *config.AppConfig
}

func NewService(ctx *common.AppContext) Service {
	return &service{
		repo:     NewRepository(ctx),
		authRepo: auth.NewRepository(ctx),
		cfg:      ctx.Config,
	}
}

func (s *service) GetAllUsers() ([]user.User, error) {
	return s.repo.GetUsers()
}

func (s *service) CreateUser(req CreateUserRequest) error {
	exists, field, err := s.authRepo.CheckUserExists(req.Username, req.Email)
	if err != nil {
		return err
	}

	if exists {
		if field == "username" {
			return auth.ErrUsernameExists
		}
		return auth.ErrEmailExists
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	newUser := &user.User{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		Profile: user.Profile{
			ID: uuid.NewString(),
			FirstName: req.FirstName,
			LastName: req.LastName,
			Gender: req.Gender,
			DOB: req.DOB,
			PhoneNumber: req.PhoneNumber,
		},
	}

	if err := s.authRepo.CreateUser(newUser); err != nil {
		return fmt.Errorf("không thể tạo người dùng: %v", err)
	}

	return nil
}
