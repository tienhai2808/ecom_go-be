package service

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/response"
)

type ProductService interface {
	GetAllProducts(ctx context.Context, query request.ProductPaginationQuery) ([]*model.Product, *response.MetaResponse, error)

	GetProductByID(ctx context.Context, id int64) (*model.Product, error)

	CreateProduct(ctx context.Context, req *request.CreateProductForm) (*model.Product, error)

	UpdateProduct(ctx context.Context, id int64, req *request.UpdateProductForm) (*model.Product, error)

	DeleteProduct(ctx context.Context, id int64) error

	DeleteProducts(ctx context.Context, req request.DeleteManyRequest) (int64, error)
}
