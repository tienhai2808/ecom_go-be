package repository

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context) ([]*model.Product, error)

	GetProductByID(ctx context.Context, id string) (*model.Product, error)

	CreateProduct(ctx context.Context, product *model.Product) error

	UpdateProductByID(ctx context.Context, id string, updateData map[string]interface{}) error

	DeleteProductByID(ctx context.Context, id string) error

	DeleteManyProducts(ctx context.Context, ids []string) (int64, error)
}
