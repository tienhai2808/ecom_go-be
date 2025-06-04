package repository

import (
	"backend/internal/model"
	"context"
)

type UserRepository interface {
	CheckUserExistsByEmail(ctx context.Context, email string) (bool, error)

	CheckUserExistsByUsername(ctx context.Context, username string) (bool, error)

	CreateUser(ctx context.Context, user *model.User) error

	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	GetUserByID(ctx context.Context, id string) (*model.User, error)

	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	
	UpdateUserPasswordByID(ctx context.Context, id, newPassword string) error
}
