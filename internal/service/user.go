package service

import (
	"backend/internal/model"
	"backend/internal/request"
	"context"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)

	GetUserByID(ctx context.Context, id string) (*model.User, error)

	CreateUser(ctx context.Context, req request.CreateUserRequest) (*model.User, error)

	UpdateUser(ctx context.Context, id string, req *request.UpdateUserRequest) (*model.User, error)

	DeleteUserByID(ctx context.Context, id string) error

	DeleteManyUsers(ctx context.Context, currentUserID string, req request.DeleteManyUsersRequest) (int64, error)
}