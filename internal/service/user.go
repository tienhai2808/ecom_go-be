package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)

	GetUserByID(ctx context.Context, id int64) (*model.User, error)

	CreateUser(ctx context.Context, req request.CreateUserRequest) (*model.User, error)

	UpdateUser(ctx context.Context, id int64, req *request.UpdateUserRequest) (*model.User, error)

	DeleteUser(ctx context.Context, id int64) error

	DeleteUsers(ctx context.Context, currentUserID int64, req request.DeleteManyRequest) (int64, error)
}
