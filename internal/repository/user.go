package repository

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"context"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]*model.User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	ExistsByUsername(ctx context.Context, username string) (bool, error)

	ExistsByID(ctx context.Context, id int64) (bool, error)

	Create(ctx context.Context, user *model.User) error

	FindByUsername(ctx context.Context, username string) (*model.User, error)

	FindByIDWithProfile(ctx context.Context, id int64) (*model.User, error)

	FindByEmail(ctx context.Context, email string) (*model.User, error)

	Update(ctx context.Context, id int64, updateData map[string]any) error

	Delete(ctx context.Context, id int64) error

	DeleteAllByID(ctx context.Context, ids []int64) (int64, error)
}
