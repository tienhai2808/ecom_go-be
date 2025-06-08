package repository

import (
	"backend/internal/model"
	"context"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]*model.User, error)

	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)

	CheckUserExistsByUsername(ctx context.Context, username string) (bool, error)

	CheckUserExistsByID(ctx context.Context, id string) (bool, error)

	CreateUser(ctx context.Context, user *model.User) error

	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	GetUserByID(ctx context.Context, id string) (*model.User, error)

	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	UpdateUserPasswordByID(ctx context.Context, id, newPassword string) error

	UpdateUserByID(ctx context.Context, id string, updateData map[string]interface{}) error

	DeleteUserByID(ctx context.Context, id string) error

	DeleteManyUsers(ctx context.Context, ids []string) (int64, error)
}
