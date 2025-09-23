package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]*model.Product, error)

	GetProductByID(ctx context.Context, id string) (*model.Product, error)

	CreateProduct(ctx context.Context, req request.CreateProductRequest) (*model.Product, error)

	UpdateProduct(ctx context.Context, id string, req *request.UpdateProductRequest) (*model.Product, error)

	DeleteProduct(ctx context.Context, id string) error

	DeleteManyProducts(ctx context.Context, req request.DeleteManyRequest) (int64, error)
}
