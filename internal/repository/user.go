package repository

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"context"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]*model.User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	ExistsByUsername(ctx context.Context, username string) (bool, error)

	ExistsByID(ctx context.Context, id string) (bool, error)

	Create(ctx context.Context, user *model.User) error

	FindByUsername(ctx context.Context, username string) (*model.User, error)

	FindByID(ctx context.Context, id string) (*model.User, error)

	FindByEmail(ctx context.Context, email string) (*model.User, error)

	UpdateUserPasswordByID(ctx context.Context, id, newPassword string) error

	Update(ctx context.Context, id string, updateData map[string]any) error

	Delete(ctx context.Context, id string) error

	DeleteAllByID(ctx context.Context, ids []string) (int64, error)
}
