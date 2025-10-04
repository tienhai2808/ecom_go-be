package service

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type ProfileService interface {
	UpdateProfile(ctx context.Context, id int64, req *request.UpdateProfileRequest) (*model.User, error)
}