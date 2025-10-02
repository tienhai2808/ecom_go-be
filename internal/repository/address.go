package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"gorm.io/gorm"
)

type AddressRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Address, error)

	FindByIDTx(ctx context.Context, tx *gorm.DB, id int64) (*model.Address, error)

	Create(ctx context.Context, address *model.Address) error

	CreateTx(ctx context.Context, tx *gorm.DB, address *model.Address) error

	FindByUserID(ctx context.Context, userID int64) ([]*model.Address, error)

	CountByUserIDTx(ctx context.Context, tx *gorm.DB, userID int64) (int64, error)

	FindLatestByUserIDExcludeIDTx(ctx context.Context, tx *gorm.DB, userID, id int64) (*model.Address, error)

	UpdateTx(ctx context.Context, tx *gorm.DB, id int64, updateData map[string]any) error

	DeleteTx(ctx context.Context, tx *gorm.DB, id int64) error
}
