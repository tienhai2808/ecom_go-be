package admin

import (
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/user"
)

type Service interface {
	GetAllUsers() ([]user.User, error)
}

type service struct {
	repo Repository
	cfg  *config.AppConfig
}

func NewService(ctx *common.AppContext) Service {
	return &service{
		repo: NewRepository(ctx),
		cfg:  ctx.Config,
	}
}

func (s *service) GetAllUsers() ([]user.User, error) {
	return s.repo.GetUsers()
}
