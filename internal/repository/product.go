package repository

import (
	"backend/internal/model"
	"context"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context) ([]*model.Product, error)

	GetProductByID(ctx context.Context, id string) (*model.Product, error)

	CreateProduct(ctx context.Context, product *model.Product) error

	UpdateProductByID(ctx context.Context, id string, updateData map[string]interface{}) error
}