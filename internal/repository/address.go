package repository

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
)

type AddressRepository interface {
	FindByID(ctx context.Context, id string) (*model.Address, error)

	Create(ctx context.Context, address *model.Address) error

	FindByUserID(ctx context.Context, userID string) ([]*model.Address, error)

	CheckDefaultExistsByUserID(ctx context.Context, userID string) (bool, error)

	UpdateNonDefaultByUserID(ctx context.Context, userID string) error

	CountByUserID(ctx context.Context, userID string) (int64, error)

	FindLatestByUserIDExcludeID(ctx context.Context, userID, id string) (*model.Address, error)

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	UpdateDefault(ctx context.Context, id string) error

	Delete(ctx context.Context, id string) error
}
