package service

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]*model.Product, error)

	SearchProduct(ctx context.Context, query string) ([]*model.Product, error)

	GetProductByID(ctx context.Context, id int64) (*model.Product, error)

	CreateProduct(ctx context.Context, req *request.CreateProductForm) (*model.Product, error)

	UpdateProduct(ctx context.Context, id int64, req *request.UpdateProductForm) (*model.Product, error)

	DeleteProduct(ctx context.Context, id int64) error

	DeleteProducts(ctx context.Context, req request.DeleteManyRequest) (int64, error)
}
