package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"gorm.io/gorm"
)

type AddressRepository interface {
	FindByID(ctx context.Context, id string) (*model.Address, error)

	FindByIDTx(ctx context.Context, tx *gorm.DB, id string) (*model.Address, error)

	Create(ctx context.Context, address *model.Address) error

	CreateTx(ctx context.Context, tx *gorm.DB, address *model.Address) error

	FindByUserID(ctx context.Context, userID string) ([]*model.Address, error)

	CountByUserIDTx(ctx context.Context, tx *gorm.DB, userID string) (int64, error)

	FindLatestByUserIDExcludeIDTx(ctx context.Context, tx *gorm.DB, userID, id string) (*model.Address, error)

	UpdateTx(ctx context.Context, tx *gorm.DB, id string, updateData map[string]any) error

	DeleteTx(ctx context.Context, tx *gorm.DB, id string) error
}
