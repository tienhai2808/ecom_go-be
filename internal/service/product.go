package service

import (
	"backend/internal/model"
	"backend/internal/request"
	"context"
)

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]*model.Product, error)

	GetProductByID(ctx context.Context, id string) (*model.Product, error)

	CreateProduct(ctx context.Context, req request.CreateProductRequest) (*model.Product, error)
}