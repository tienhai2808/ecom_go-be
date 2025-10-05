package repository

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]*model.Product, error)

	FindByID(ctx context.Context, id int64) (*model.Product, error)

	Create(ctx context.Context, product *model.Product) error

	Update(ctx context.Context, id int64, updateData map[string]any) error

	Delete(ctx context.Context, id int64) error

	DeleteAllByID(ctx context.Context, ids []int64) (int64, error)
}
