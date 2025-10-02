package repository

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context) ([]*model.Product, error)

	GetProductByID(ctx context.Context, id int64) (*model.Product, error)

	CreateProduct(ctx context.Context, product *model.Product) error

	UpdateProductByID(ctx context.Context, id int64, updateData map[string]any) error

	DeleteProductByID(ctx context.Context, id int64) error

	DeleteManyProducts(ctx context.Context, ids []int64) (int64, error)
}
