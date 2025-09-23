package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)

	GetUserByID(ctx context.Context, id string) (*model.User, error)

	CreateUser(ctx context.Context, req request.CreateUserRequest) (*model.User, error)

	UpdateUser(ctx context.Context, id string, req *request.UpdateUserRequest) (*model.User, error)

	DeleteUserByID(ctx context.Context, id string) error

	DeleteManyUsers(ctx context.Context, currentUserID string, req request.DeleteManyRequest) (int64, error)
}
