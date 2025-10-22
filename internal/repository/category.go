package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error

	FindAll(ctx context.Context) ([]*model.Category, error)
	
	FindByID(ctx context.Context, id int64) (*model.Category, error)

	FindByIDTx(ctx context.Context, tx *gorm.DB, id int64) (*model.Category, error)

	Update(ctx context.Context, id int64, updateData map[string]any) error
}